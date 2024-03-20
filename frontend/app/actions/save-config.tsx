'use server'

export async function saveConfig(formData: any) {
    console.log('formData:', formData);
    const username = formData.get('username');
    const sshPublicKeyFile = formData.get('sshPublicKeyFile');
    const pullSecretFile = formData.get('pullSecretFile');
    const outputDir = formData.get('outputDir');
    const clusterName = formData.get('clusterName');
    const image = formData.get('image');
    const cloudRegion = formData.get('cloudRegion');
    const cloud = formData.get('cloud');
    const dryRun = formData.get('dryRun');

    let result : Response | any;
    try {
        let responseBody = JSON.stringify({
            "username": username,
            "sshPublicKeyFile": sshPublicKeyFile,
            "pullSecretFile": pullSecretFile,
            "outputDir": outputDir,
            "clusterName": clusterName,
            "image": image,
            "cloudRegion": cloudRegion,
            "cloud": cloud,
            "dryRun": dryRun? "true" : "false",
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