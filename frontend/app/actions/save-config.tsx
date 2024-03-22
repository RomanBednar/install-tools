'use server'

export async function saveConfig(previousState: any, formData: { get: (arg0: string) => any; }) {
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

    let result : {message: string, success: boolean};
    try {
        let requestBody = JSON.stringify({
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
        console.log('requestBody:', requestBody)
        const apiUrl = process.env.NEXT_PUBLIC_API_URL;
        console.log("Connecting to:", apiUrl)
        const response = await fetch(`${apiUrl}/save`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: requestBody,
        });
        result = {message: "OK", success: true};

    } catch (error) {
        result = {message: JSON.stringify(error), success: true};
    }
    return result
}