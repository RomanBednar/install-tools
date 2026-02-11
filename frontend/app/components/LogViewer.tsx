"use client";

import React, { useState, useEffect, useRef } from "react";

interface LogViewerProps {
  active: boolean;
  outputDir: string;
}

export default function LogViewer({ active, outputDir }: LogViewerProps) {
  const [logContent, setLogContent] = useState("");
  const logRef = useRef<HTMLPreElement>(null);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (!active) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      return;
    }

    const fetchLog = async () => {
      try {
        const logPath = outputDir
          ? `${outputDir}/.openshift_install.log`
          : "";
        const url = logPath
          ? `/api/log?path=${encodeURIComponent(logPath)}`
          : "/api/log";
        const resp = await fetch(url);
        if (resp.ok) {
          const data = await resp.json();
          if (data.log) {
            setLogContent(data.log);
          }
        }
      } catch {
        // Ignore fetch errors during polling
      }
    };

    // Initial fetch
    fetchLog();

    // Poll every 2 seconds
    intervalRef.current = setInterval(fetchLog, 2000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [active, outputDir]);

  // Auto-scroll to bottom
  useEffect(() => {
    if (logRef.current) {
      logRef.current.scrollTop = logRef.current.scrollHeight;
    }
  }, [logContent]);

  if (!active) return null;

  return (
    <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-900 dark:bg-black overflow-hidden">
      <div className="flex items-center justify-between px-4 py-2 bg-gray-800 dark:bg-gray-900 border-b border-gray-700">
        <div className="flex items-center gap-2">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-red-500" />
            <div className="w-3 h-3 rounded-full bg-yellow-500" />
            <div className="w-3 h-3 rounded-full bg-green-500" />
          </div>
          <span className="text-xs text-gray-400 ml-2 font-mono">
            .openshift_install.log
          </span>
        </div>
        <div className="flex items-center gap-1">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75" />
            <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500" />
          </span>
          <span className="text-xs text-green-400">Live</span>
        </div>
      </div>
      <pre
        ref={logRef}
        className="log-viewer p-4 text-xs text-green-400 font-mono whitespace-pre-wrap break-all overflow-y-auto"
        style={{ maxHeight: "400px", minHeight: "200px" }}
      >
        {logContent || "Waiting for log output..."}
      </pre>
    </div>
  );
}
