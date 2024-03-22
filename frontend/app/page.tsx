'use client'

import { saveConfig } from '@/app/actions/save-config';
import { SubmitButton } from './ui/submit-button';
import React, { useState } from 'react';
import DefaultInput from "@/app/ui/default-input";
import CloudSelect from "@/app/ui/cloud-selection";
import DryRunSwitch from "@/app/ui/dryrun-switch";
import ActionButton from "@/app/ui/action-button";
import FileViewer from "@/app/ui/file-viewer";


export default function InstallerForm() {

    const [inputValue, setInputValue] = useState('');

    const handleUseDefault = () => {
        setInputValue('aaa');
    };

    return (
        <div>
            <form action={saveConfig}>
                <h1 className="py-5 text-3xl font-semibold leading-9 text-gray-900">OpenShift Installer</h1>
                <div className="space-y-12 md:container md:mx-auto">
                    <div className="border-b border-gray-900/10 pb-12">
                        <h2 className="text-base font-semibold leading-7 text-gray-900">User & configuration</h2>
                        <p className="mt-1 text-sm leading-6 text-gray-600">
                            Information about user and the configuration of the installer.
                        </p>

                        <div>
                            <DefaultInput name="username" placeholder="Your username" label="Username"
                                          disableDefaultButton={true}/>
                        </div>
                        <div>
                            <DefaultInput name="sshPublicKeyFile" defaultValue="${HOME}/.ssh/id_rsa.pub"
                                          placeholder="${HOME}/.ssh/id_rsa.pub" label="Public SSH key"/>
                        </div>
                        <div>
                            <DefaultInput name="pullSecretFile" defaultValue="${HOME}/.docker/config.json"
                                          placeholder="${HOME}/.docker/config.json" label="Pull secret file"/>
                        </div>
                        <div>
                            <DefaultInput name="outputDir" defaultValue="/tmp/output"
                                          placeholder="/tmp/output" label="Output directory"/>
                            <p className="text-sm text-red-600">Warning: anything outside home dir (readonly) is stored
                                in a
                                container.</p>
                            <p className="text-sm text-red-600">All data is lost if application is stopped.</p>
                        </div>
                    </div>

                    <div className="border-b border-gray-900/10 pb-12">
                        <h2 className="text-base font-semibold leading-7 text-gray-900">Cluster Information</h2>
                        <p className="mt-1 text-sm leading-6 text-gray-600">Details about cluster environment like cloud
                            platform, payload image or region.</p>

                        <div>
                            <DefaultInput name="clusterName" placeholder="Name of the cluster" label="Cluster name"
                                          disableDefaultButton={true}/>
                        </div>

                        <div>
                            <DefaultInput name="image" placeholder="OpenShift image to use" label="Payload image"
                                          disableDefaultButton={true}/>
                            <p className="text-sm">Hint: select image <a className="text-indigo-600" target="_blank"
                                                                         href="https://amd64.ocp.releases.ci.openshift.org/">here</a>
                            </p>
                        </div>

                        <div>
                            <DefaultInput name="cloudRegion" placeholder="eu-central-1" label="Region"
                                          disableDefaultButton={true}/>
                        </div>

                        <div>
                            <CloudSelect/>
                        </div>
                        <div>
                            <DryRunSwitch/>
                        </div>
                        <div className="mt-6 flex items-center justify-start gap-x-6">
                            <SubmitButton/>
                        </div>
                    </div>


                </div>

            </form>

            <div className="space-y-12 md:container md:mx-auto">
                <div className="border-b border-gray-900/10 pb-12">
                    <h2 className="text-base font-semibold leading-7 text-gray-900">Actions: </h2>
                    <ActionButton action="create"/>
                    <ActionButton action="destroy"/>
                </div>
            </div>

            <div className="space-y-12 md:container md:mx-auto">
                <div className="border-b border-gray-900/10 pb-12">
                    <h2 className="text-base font-semibold leading-7 text-gray-900">Progress: </h2>

                </div>
            </div>
        </div>

    )
}