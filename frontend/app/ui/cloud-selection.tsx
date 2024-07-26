import React from "react";


function CloudSelect() {

    const cloudPlatforms = ["aws", "aws-odf", "azure"]

    return (
        <div className="py-5">
            <fieldset>
                <legend className="text-sm font-semibold leading-6 text-gray-900">Cloud Platform
                </legend>
                <p className="mt-1 text-sm leading-6 text-gray-600">Select a cloud platform or variant.</p>
                <div className="mt-6 space-y-6">
                    {cloudPlatforms.map(cloud =>

                        <div key={cloud} className="flex items-center gap-x-3">
                            <input
                                id={cloud}
                                value={cloud}
                                name="cloud"
                                type="radio"
                                className="h-4 w-4 border-gray-300 text-indigo-600 focus:ring-indigo-600"
                            />
                            <label htmlFor={cloud}
                                   className="block text-sm font-medium leading-6 text-gray-900">
                                {cloud}
                            </label>
                        </div>
                    )}
                </div>
            </fieldset>
        </div>


    )
}

export default CloudSelect;