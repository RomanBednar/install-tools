'use server'

export async function saveConfig(formData: any) {
    console.log('formData:', formData);
    const user = formData.get('username');
    const pwd = formData.get('password');
    // const response = await fetch(`https://api.example.com/articles/${articleId}/comments`, {
    //     method: 'POST',
    //     headers: {
    //         'Content-Type': 'application/json',
    //     },
    //     body: JSON.stringify({ comment }),
    // });
    let result : Response | any;
    try {
        let responseBody = JSON.stringify({
            "username": user,
            "password": pwd
        })
        console.log('responseBody:', responseBody)
        const apiUrl = process.env.NEXT_PUBLIC_API_URL;
        console.log("Connecting to:", apiUrl)
        const response = await fetch(`${apiUrl}/storeCredentials`, {
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