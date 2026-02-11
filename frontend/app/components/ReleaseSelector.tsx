"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";

interface Release {
  name: string;
  phase: string;
  imageUrl: string;
}

interface ReleaseSelectorProps {
  value: string;
  onChange: (imageUrl: string, releaseName: string) => void;
}

export default function ReleaseSelector({
  value,
  onChange,
}: ReleaseSelectorProps) {
  const [releases, setReleases] = useState<Release[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [selectedName, setSelectedName] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const fetchReleases = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const resp = await fetch("/api/releases");
      if (!resp.ok) throw new Error("Failed to fetch releases");
      const data = await resp.json();
      setReleases(data.releases || []);
    } catch (e) {
      setError("Failed to load releases. Please try again.");
      console.error("Error fetching releases:", e);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchReleases();
  }, [fetchReleases]);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const filtered = releases.filter((r) =>
    r.name.toLowerCase().includes(query.toLowerCase())
  );

  const handleSelect = (release: Release) => {
    setSelectedName(release.name);
    setQuery(release.name);
    setIsOpen(false);
    onChange(release.imageUrl, release.name);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setQuery(val);
    setIsOpen(true);
    // If user clears the input, clear the selection
    if (!val) {
      setSelectedName("");
      onChange("", "");
    }
  };

  return (
    <div className="relative">
      <div className="flex items-center gap-2">
        <div className="relative flex-1">
          <input
            ref={inputRef}
            type="text"
            value={query || selectedName}
            onChange={handleInputChange}
            onFocus={() => setIsOpen(true)}
            placeholder={loading ? "Loading releases..." : "Type to search releases (e.g. 4.22)"}
            className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-none transition-colors"
            disabled={loading}
          />
          {selectedName && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2">
              <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 dark:bg-green-900/50 text-green-800 dark:text-green-300">
                Accepted
              </span>
            </div>
          )}
        </div>
        <button
          type="button"
          onClick={fetchReleases}
          disabled={loading}
          className="inline-flex items-center justify-center w-10 h-10 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700 hover:text-indigo-600 dark:hover:text-indigo-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-colors disabled:opacity-50"
          title="Refresh releases"
        >
          <svg
            className={`w-4 h-4 ${loading ? "animate-spin" : ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        </button>
      </div>

      {error && (
        <p className="mt-1 text-xs text-red-500 dark:text-red-400">{error}</p>
      )}

      {isOpen && !loading && filtered.length > 0 && (
        <div
          ref={dropdownRef}
          className="absolute z-40 mt-1 w-full max-h-60 overflow-y-auto rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-lg"
        >
          {filtered.map((release) => (
            <button
              key={release.name}
              type="button"
              onClick={() => handleSelect(release)}
              className={`w-full text-left px-4 py-2 text-sm hover:bg-indigo-50 dark:hover:bg-indigo-900/30 transition-colors ${
                selectedName === release.name
                  ? "bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300"
                  : "text-gray-700 dark:text-gray-300"
              }`}
            >
              <span className="font-mono">{release.name}</span>
              <span className="ml-2 inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-green-100 dark:bg-green-900/50 text-green-700 dark:text-green-300">
                Accepted
              </span>
            </button>
          ))}
        </div>
      )}

      {isOpen && !loading && query && filtered.length === 0 && (
        <div className="absolute z-40 mt-1 w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-lg p-4 text-center text-sm text-gray-500 dark:text-gray-400">
          No matching releases found
        </div>
      )}

      {value && (
        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400 font-mono truncate">
          Image: {value}
        </p>
      )}
    </div>
  );
}
