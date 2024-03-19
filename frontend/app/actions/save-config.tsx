'use server'

export async function saveConfig(formData: any) {
    console.log('formData:', formData);
    const username = formData.get('username');
    const sshPublicKeyFile = formData.get('sshPublicKeyFile');

    let result : Response | any;
    try {
        let responseBody = JSON.stringify({
            "username": username,
            "sshPublicKeyFile": sshPublicKeyFile
        })
        console.log('responseBody:', responseBody)
        const apiUrl = process.env.NEXT_PUBLIC_API_URL;
        console.log("Connecting to:", apiUrl)
        const response = await fetch(`${apiUrl}/save`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: responseBody,
        });
        result = response;
    } catch (error) {
        console.error('Error:', error);
        result = error;
    }

    console.log('Result:', result);

    return JSON.stringify(result);
}