export const dynamic = 'force-dynamic';

async function getData() {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL;
    console.log("Connecting to:", apiUrl)
    const res = await fetch(`${apiUrl}/log`)
    if (!res.ok) {
        throw new Error('Failed to fetch data')
    }
    return res.text()
}

export default async function Page() {
    const data = await getData()

    return (
        <div>
            <p>{data}</p>
        </div>

    )
}