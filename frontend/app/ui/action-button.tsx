// components/Button.tsx
import React from 'react';
import { runAction } from '@/app/actions/run-action';
import clsx from 'clsx';


interface ButtonProps {
    action: 'create' | 'destroy';
}

const ActionButton: React.FC<ButtonProps> = ({ action }) => {
    const handleClick = async () => {
        try {
            const result = await runAction(action);
            console.log('Result:', result);
            // Handle result as needed
        } catch (error) {
            console.error('Error:', error);
            // Handle error
        }
    };

    return (
        <button onClick={handleClick} className={clsx(
        "rounded-md mr-5 px-3 py-2 text-sm font-semibold text-white shadow-sm  focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
            {
                'bg-green-600 hover:bg-green-500': action === 'create',
                'bg-red-600 hover:bg-red-500': action === 'destroy',
            },
            )}
        >
            {action}
        </button>
    );
};

export default ActionButton;
