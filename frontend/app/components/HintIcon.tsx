"use client";

import React, { useState } from "react";

interface HintIconProps {
  text: string;
}

export default function HintIcon({ text }: HintIconProps) {
  const [show, setShow] = useState(false);

  return (
    <span className="relative inline-block ml-1.5">
      <button
        type="button"
        className="inline-flex items-center justify-center w-4 h-4 rounded-full bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-200 text-xs font-bold hover:bg-indigo-400 hover:text-white dark:hover:bg-indigo-500 transition-colors cursor-help"
        onMouseEnter={() => setShow(true)}
        onMouseLeave={() => setShow(false)}
        onClick={() => setShow(!show)}
        aria-label="Info"
      >
        i
      </button>
      {show && (
        <div className="absolute z-50 bottom-full left-1/2 -translate-x-1/2 mb-2 w-72 p-3 rounded-lg shadow-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-xs text-gray-700 dark:text-gray-300 leading-relaxed">
          {text}
          <div className="absolute top-full left-1/2 -translate-x-1/2 w-2 h-2 bg-white dark:bg-gray-800 border-b border-r border-gray-200 dark:border-gray-700 rotate-45 -mt-1" />
        </div>
      )}
    </span>
  );
}
