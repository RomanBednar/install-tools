import { useEffect, useState } from 'react';

interface FileViewerProps {
    filePath: string;
}

const FileViewer: React.FC<FileViewerProps> = ({ filePath }) => {
    const [fileContent, setFileContent] = useState<string[]>([]);

    useEffect(() => {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL;
        console.log("Connecting to:", apiUrl)
        const eventSource = new EventSource(`${apiUrl}/log`);

        eventSource.onmessage = (event) => {
            setFileContent(prevContent => [...prevContent, event.data]);
        };

        return () => {
            eventSource.close();
        };
    }, [filePath]);

    return (
        <div
            className="file-viewer"
            style={{ width: '500px', height: '300px', overflow: 'auto', border: '1px solid #ccc' }}
        >
            {fileContent.map((line, index) => (
                <div key={index}>{line}</div>
            ))}
        </div>
    );
};

export default FileViewer;