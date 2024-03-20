'use server'

export async function runAction(action: string) {

    let result : Response | any;
    try {
        let responseBody = JSON.stringify({action});
        console.log('responseBody:', responseBody)
        const apiUrl = process.env.NEXT_PUBLIC_API_URL;
        console.log("Connecting to:", apiUrl)
        const response = await fetch(`${apiUrl}/action`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: responseBody,
        });
        result = response;
        console.error('Response:', result);
    } catch (error) {
        console.error('Error:', error);
        result = error;
    }

    console.log('Result:', result);

    return JSON.stringify(result);
}