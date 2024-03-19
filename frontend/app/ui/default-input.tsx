import React, { useState } from 'react';

export interface InputInterface {
    name: string;
    placeholder: string;
    label: string;
    defaultValue?: string;
    disableDefaultButton?: boolean;
}

function DefaultInput({name, defaultValue, placeholder, label, disableDefaultButton} : InputInterface) {

    const [inputValue, setInputValue] = useState('');

    const handleUseDefault = () => {
        setInputValue(defaultValue? defaultValue : "");
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value);
    };

    return (
        <div className="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
            <div className="sm:col-span-4">
                <label htmlFor={name} className="block text-sm font-medium leading-6 text-gray-900">
                    {label}
                </label>
                <div className="mt-2">
                    <div
                        className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
                        <input
                            type="text"
                            name={name}
                            id={name}
                            className="block flex-1 border-0 bg-transparent py-1.5 pl-1 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
                            placeholder={placeholder}
                            value={inputValue ? inputValue : ""}
                            onChange={handleChange}
                        />
                        {!disableDefaultButton &&
                        <button onClick={handleUseDefault}
                                className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                            Use default
                        </button>
                        }

                    </div>
                </div>
            </div>

        </div>
    );
}

export default DefaultInput;