import React from "react";


function DryRunSwitch() {
    return (
        <div className="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
            <div className="sm:col-span-4">
                <legend className="text-sm font-semibold leading-6 text-gray-900">Dry run
                </legend>
                <p className="mt-1 text-sm leading-6 text-gray-600">Only prepare configuration files and don't run openshift-install.</p>
                <div className="m-2 flex">
                    <label htmlFor="dryRun" className="block text-sm font-medium leading-6 text-gray-900 relative">
                        Enable
                    </label>
                    <input
                        type="checkbox"
                        name="dryRun"
                        id="dryRun"
                        className="w-4 ml-2"
                    />
                </div>
            </div>
        </div>
    );
}

export default DryRunSwitch;