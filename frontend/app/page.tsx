"use client";

import React, { useState, useEffect, useCallback } from "react";
import ThemeToggle from "@/app/components/ThemeToggle";
import HintIcon from "@/app/components/HintIcon";
import ReleaseSelector from "@/app/components/ReleaseSelector";
import LogViewer from "@/app/components/LogViewer";

const CLOUD_PLATFORMS = [
  { value: "aws", label: "AWS" },
  { value: "aws-sts", label: "AWS STS" },
  { value: "aws-odf", label: "AWS ODF" },
  { value: "azure", label: "Azure" },
  { value: "azure-wi", label: "Azure Workload Identity" },
  { value: "gcp", label: "GCP" },
  { value: "gcp-wif", label: "GCP WIF" },
  { value: "vsphere", label: "vSphere" },
  { value: "alibaba", label: "Alibaba" },
];

const ACTIONS = [
  { value: "create", label: "Create" },
  { value: "destroy", label: "Destroy" },
];

type OperationStatus = "idle" | "running" | "completed" | "error" | "unknown";

export default function InstallerPage() {
  // Theme
  const [isDark, setIsDark] = useState(false);

  // Form state
  const [username, setUsername] = useState("");
  const [sshPublicKeyFile, setSshPublicKeyFile] = useState("");
  const [pullSecretFile, setPullSecretFile] = useState("");
  const [outputDir, setOutputDir] = useState("");
  const [clusterName, setClusterName] = useState("");
  const [image, setImage] = useState("");
  const [releaseName, setReleaseName] = useState("");
  const [cloudRegion, setCloudRegion] = useState("");
  const [cloud, setCloud] = useState("aws");
  const [action, setAction] = useState("create");
  const [dryRun, setDryRun] = useState(false);

  // Operation state
  const [operationStatus, setOperationStatus] =
    useState<OperationStatus>("idle");
  const [operationError, setOperationError] = useState("");
  const [operationMessage, setOperationMessage] = useState("");
  const [dirEmpty, setDirEmpty] = useState(true);
  const [dirChecking, setDirChecking] = useState(false);
  const [dirWarning, setDirWarning] = useState("");
  const [formError, setFormError] = useState("");

  // Completed dry-run message
  const [dryRunResult, setDryRunResult] = useState("");

  // Initialize theme from localStorage
  useEffect(() => {
    const stored = localStorage.getItem("theme");
    if (stored === "dark" || (!stored && window.matchMedia("(prefers-color-scheme: dark)").matches)) {
      setIsDark(true);
      document.documentElement.classList.add("dark");
    }
  }, []);

  const toggleTheme = () => {
    const next = !isDark;
    setIsDark(next);
    if (next) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }
  };

  // Check output directory when it changes
  useEffect(() => {
    if (!outputDir.trim()) {
      setDirWarning("");
      setDirEmpty(true);
      return;
    }

    const checkDir = async () => {
      setDirChecking(true);
      try {
        const resp = await fetch(
          `/api/check-dir?path=${encodeURIComponent(outputDir)}`
        );
        if (resp.ok) {
          const data = await resp.json();
          if (data.exists && !data.empty) {
            setDirEmpty(false);
            setDirWarning(
              "This directory is not empty. Please choose an empty directory or clear its contents."
            );
          } else {
            setDirEmpty(true);
            setDirWarning("");
          }
        }
      } catch {
        // If we can't check, allow it
        setDirEmpty(true);
        setDirWarning("");
      } finally {
        setDirChecking(false);
      }
    };

    const timer = setTimeout(checkDir, 500);
    return () => clearTimeout(timer);
  }, [outputDir]);

  // Poll operation status when running
  useEffect(() => {
    if (operationStatus !== "running") return;

    const pollStatus = async () => {
      try {
        const resp = await fetch("/api/status");
        if (resp.ok) {
          const data = await resp.json();
          setOperationStatus(data.status as OperationStatus);
          setOperationError(data.error || "");
          setOperationMessage(data.message || "");

          if (data.status === "completed" && dryRun) {
            // Convert container path to host path for user info
            const containerPrefix = "/root/ocp-install-tool/";
            const hostPrefix = "~/ocp-install-tool/";
            let displayPath = outputDir;
            if (outputDir.startsWith(containerPrefix)) {
              displayPath = outputDir.replace(containerPrefix, hostPrefix);
            } else if (outputDir.startsWith("/root/")) {
              displayPath = outputDir.replace("/root/", "~/");
            }
            setDryRunResult(
              `Dry run completed successfully! Installation has been prepared in: ${displayPath}`
            );
          }
        }
      } catch {
        // ignore polling errors
      }
    };

    const interval = setInterval(pollStatus, 2000);
    return () => clearInterval(interval);
  }, [operationStatus, dryRun, outputDir]);

  // When output dir changes, reset operation status so user can start new operation
  useEffect(() => {
    if (operationStatus === "completed" || operationStatus === "error") {
      setOperationStatus("idle");
      setOperationError("");
      setOperationMessage("");
      setDryRunResult("");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [outputDir]);

  const handleStart = useCallback(async () => {
    setFormError("");
    setDryRunResult("");

    // Validate required fields
    if (!clusterName.trim()) {
      setFormError("Cluster name is required.");
      return;
    }
    if (!image.trim()) {
      setFormError("Release image is required. Please select a release.");
      return;
    }
    if (!outputDir.trim()) {
      setFormError("Output directory is required.");
      return;
    }
    if (!dirEmpty) {
      setFormError(
        "Output directory must be empty. Please choose a different directory."
      );
      return;
    }

    setOperationStatus("running");
    setOperationError("");
    setOperationMessage("Starting operation...");

    try {
      const resp = await fetch("/api/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username,
          sshPublicKeyFile,
          pullSecretFile,
          outputDir,
          clusterName,
          image,
          cloudRegion,
          cloud,
          action,
          dryRun,
        }),
      });

      if (!resp.ok) {
        const data = await resp.json();
        setOperationStatus("error");
        setOperationError(data.error || "Failed to start operation");
        return;
      }

      // Operation started successfully - status polling will track it
    } catch (error) {
      setOperationStatus("error");
      setOperationError(`Failed to connect: ${error}`);
    }
  }, [
    username,
    sshPublicKeyFile,
    pullSecretFile,
    outputDir,
    clusterName,
    image,
    cloudRegion,
    cloud,
    action,
    dryRun,
    dirEmpty,
  ]);

  const isStartDisabled =
    operationStatus === "running" ||
    operationStatus === "completed" ||
    !dirEmpty ||
    dirChecking;

  // Show log viewer: when NOT dry-run and operation is running or completed
  const showLogViewer = !dryRun && (operationStatus === "running" || operationStatus === "completed" || operationStatus === "error");

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950 transition-colors duration-300">
      {/* Header */}
      <header className="sticky top-0 z-30 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-gradient-to-br from-red-500 to-red-700 rounded-lg flex items-center justify-center">
              <svg
                className="w-5 h-5 text-white"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
              </svg>
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900 dark:text-white">
                OpenShift Installer
              </h1>
              <p className="text-xs text-gray-500 dark:text-gray-400">
                Cluster deployment GUI
              </p>
            </div>
          </div>
          <ThemeToggle isDark={isDark} onToggle={toggleTheme} />
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Status Bar */}
        {operationStatus !== "idle" && (
          <div
            className={`mb-6 rounded-xl p-4 border ${
              operationStatus === "running"
                ? "bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800"
                : operationStatus === "completed"
                ? "bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800"
                : operationStatus === "error"
                ? "bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800"
                : "bg-gray-50 dark:bg-gray-800 border-gray-200 dark:border-gray-700"
            }`}
          >
            <div className="flex items-center gap-3">
              {operationStatus === "running" && (
                <svg
                  className="animate-spin h-5 w-5 text-blue-600 dark:text-blue-400"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
              )}
              {operationStatus === "completed" && (
                <svg
                  className="h-5 w-5 text-green-600 dark:text-green-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              )}
              {operationStatus === "error" && (
                <svg
                  className="h-5 w-5 text-red-600 dark:text-red-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              )}
              <div>
                <p
                  className={`text-sm font-medium ${
                    operationStatus === "running"
                      ? "text-blue-800 dark:text-blue-300"
                      : operationStatus === "completed"
                      ? "text-green-800 dark:text-green-300"
                      : operationStatus === "error"
                      ? "text-red-800 dark:text-red-300"
                      : "text-gray-800 dark:text-gray-300"
                  }`}
                >
                  {operationMessage || `Status: ${operationStatus}`}
                </p>
                {operationError && (
                  <p className="text-xs text-red-600 dark:text-red-400 mt-1">
                    {operationError}
                  </p>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Dry Run Result */}
        {dryRunResult && (
          <div className="mb-6 rounded-xl p-4 border bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800">
            <div className="flex items-start gap-3">
              <svg
                className="h-5 w-5 text-green-600 dark:text-green-400 mt-0.5 flex-shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <p className="text-sm text-green-800 dark:text-green-300">
                {dryRunResult}
              </p>
            </div>
          </div>
        )}

        <div className="space-y-8">
          {/* Section: User & Configuration */}
          <section className="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-800 overflow-hidden">
            <div className="px-6 py-5 border-b border-gray-100 dark:border-gray-800">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                User & Configuration
              </h2>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Information about user and the configuration of the installer.
              </p>
            </div>
            <div className="px-6 py-6 space-y-5">
              {/* Username */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  Username
                  <HintIcon text="The username for the cluster installation. Used to identify resources created during installation." />
                </label>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="mytestuser-1"
                  className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-none transition-colors"
                />
              </div>

              {/* SSH Public Key */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  SSH Public Key File
                  <HintIcon text="Path to the SSH public key file that will be injected into cluster nodes. Inside the container, your home directory is mounted at /root." />
                </label>
                <input
                  type="text"
                  value={sshPublicKeyFile}
                  onChange={(e) => setSshPublicKeyFile(e.target.value)}
                  placeholder="/root/.ssh/id_rsa.pub"
                  className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-none transition-colors"
                />
              </div>

              {/* Pull Secret */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  Pull Secret File
                  <HintIcon text="Path to the pull secret file for authenticating with container registries. Inside the container, your home directory is mounted at /root." />
                </label>
                <input
                  type="text"
                  value={pullSecretFile}
                  onChange={(e) => setPullSecretFile(e.target.value)}
                  placeholder="/root/.docker/config.json"
                  className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-none transition-colors"
                />
              </div>

              {/* Output Directory */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  Output Directory
                  <HintIcon text="The output directory inside the container. Your home directory is mounted under /root, so any path like /root/ocp-install-tool/clusters/aws/my-cluster will be available at ~/ocp-install-tool/clusters/aws/my-cluster on your host machine." />
                </label>
                <input
                  type="text"
                  value={outputDir}
                  onChange={(e) => setOutputDir(e.target.value)}
                  placeholder="/root/ocp-install-tool/clusters/aws/my-cluster"
                  className={`w-full rounded-lg border px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:ring-2 focus:outline-none transition-colors bg-white dark:bg-gray-800 ${
                    dirWarning
                      ? "border-red-400 dark:border-red-500 focus:border-red-500 focus:ring-red-500/20"
                      : "border-gray-300 dark:border-gray-600 focus:border-indigo-500 focus:ring-indigo-500/20"
                  }`}
                />
                {dirWarning && (
                  <p className="mt-1.5 text-xs text-red-600 dark:text-red-400 flex items-center gap-1">
                    <svg
                      className="w-3.5 h-3.5 flex-shrink-0"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                        clipRule="evenodd"
                      />
                    </svg>
                    {dirWarning}
                  </p>
                )}
              </div>
            </div>
          </section>

          {/* Section: Cluster Information */}
          <section className="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-800 overflow-hidden">
            <div className="px-6 py-5 border-b border-gray-100 dark:border-gray-800">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Cluster Information
              </h2>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                Details about the cluster environment like cloud platform,
                payload image, or region.
              </p>
            </div>
            <div className="px-6 py-6 space-y-5">
              {/* Cluster Name */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  Cluster Name
                  <HintIcon text="Name of the OpenShift cluster to create. This will be used as a prefix for all cloud resources." />
                </label>
                <input
                  type="text"
                  value={clusterName}
                  onChange={(e) => setClusterName(e.target.value)}
                  placeholder="mytestcluster-1"
                  className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-none transition-colors"
                />
              </div>

              {/* Release Image */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  Release Image
                  <HintIcon text="Select an OpenShift release image. Only 'Accepted' (green) releases from the official release page are shown. Type to filter versions (e.g., type '4.22' to see 4.22.x releases)." />
                </label>
                <ReleaseSelector
                  value={image}
                  onChange={(imageUrl, name) => {
                    setImage(imageUrl);
                    setReleaseName(name);
                  }}
                />
              </div>

              {/* Cloud Region */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  Cloud Region
                  <HintIcon text="The cloud provider region where the cluster will be deployed. Examples: us-east-1 (AWS), eastus (Azure), us-central1 (GCP)." />
                </label>
                <input
                  type="text"
                  value={cloudRegion}
                  onChange={(e) => setCloudRegion(e.target.value)}
                  placeholder="us-east-1"
                  className="w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-2.5 text-sm text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-none transition-colors"
                />
              </div>

              {/* Cloud Platform */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Cloud Platform
                  <HintIcon text="The cloud provider and variant to use for installation. Choose the platform that matches your cloud account and desired deployment method." />
                </label>
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                  {CLOUD_PLATFORMS.map((cp) => (
                    <button
                      key={cp.value}
                      type="button"
                      onClick={() => setCloud(cp.value)}
                      className={`px-4 py-2.5 rounded-lg border text-sm font-medium transition-all ${
                        cloud === cp.value
                          ? "bg-indigo-600 border-indigo-600 text-white shadow-md shadow-indigo-500/25"
                          : "bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:border-indigo-300 dark:hover:border-indigo-600 hover:bg-indigo-50 dark:hover:bg-indigo-900/20"
                      }`}
                    >
                      {cp.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Action */}
              <div>
                <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Action
                  <HintIcon text="Choose 'Create' to deploy a new cluster or 'Destroy' to tear down an existing one." />
                </label>
                <div className="flex gap-3">
                  {ACTIONS.map((a) => (
                    <button
                      key={a.value}
                      type="button"
                      onClick={() => setAction(a.value)}
                      className={`flex-1 px-4 py-2.5 rounded-lg border text-sm font-medium transition-all ${
                        action === a.value
                          ? a.value === "create"
                            ? "bg-green-600 border-green-600 text-white shadow-md shadow-green-500/25"
                            : "bg-red-600 border-red-600 text-white shadow-md shadow-red-500/25"
                          : "bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-500"
                      }`}
                    >
                      {a.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Dry Run Toggle */}
              <div className="flex items-center justify-between rounded-lg border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 p-4">
                <div>
                  <label className="flex items-center text-sm font-medium text-gray-700 dark:text-gray-300">
                    Dry Run
                    <HintIcon text="When enabled, only configuration files (install-config.yaml) and tools (oc, openshift-install) will be prepared. No actual cluster installation will happen." />
                  </label>
                  <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                    Generate configs only, do not install the cluster
                  </p>
                </div>
                <button
                  type="button"
                  role="switch"
                  aria-checked={dryRun}
                  onClick={() => setDryRun(!dryRun)}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    dryRun
                      ? "bg-indigo-600"
                      : "bg-gray-300 dark:bg-gray-600"
                  }`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform ${
                      dryRun ? "translate-x-6" : "translate-x-1"
                    }`}
                  />
                </button>
              </div>
            </div>
          </section>

          {/* Start Button */}
          <div className="flex flex-col items-center gap-3">
            {formError && (
              <p className="text-sm text-red-600 dark:text-red-400 text-center">
                {formError}
              </p>
            )}
            <button
              type="button"
              onClick={handleStart}
              disabled={isStartDisabled}
              className={`w-full sm:w-auto px-8 py-3.5 rounded-xl text-base font-semibold shadow-lg transition-all ${
                isStartDisabled
                  ? "bg-gray-300 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed shadow-none"
                  : "bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-500 hover:to-indigo-600 text-white shadow-indigo-500/30 hover:shadow-indigo-500/40 hover:-translate-y-0.5 active:translate-y-0"
              }`}
            >
              {operationStatus === "running" ? (
                <span className="flex items-center gap-2">
                  <svg
                    className="animate-spin h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                  </svg>
                  Running...
                </span>
              ) : operationStatus === "completed" ? (
                "Completed"
              ) : (
                `Start ${action === "create" ? "Installation" : "Destroy"}`
              )}
            </button>
            {operationStatus === "completed" && (
              <p className="text-xs text-gray-500 dark:text-gray-400">
                Change the output directory to start a new operation.
              </p>
            )}
          </div>

          {/* Live Log Viewer */}
          {showLogViewer && (
            <section>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Live Log Viewer
              </h2>
              <LogViewer active={showLogViewer} outputDir={outputDir} />
            </section>
          )}
        </div>
      </main>

      {/* Footer */}
      <footer className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6 mt-8 border-t border-gray-200 dark:border-gray-800">
        <p className="text-center text-xs text-gray-400 dark:text-gray-600">
          OpenShift Installer GUI &mdash; Built with Next.js + Tailwind CSS
        </p>
      </footer>
    </div>
  );
}
