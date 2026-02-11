import { NextResponse } from "next/server";

export const dynamic = "force-dynamic";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function GET() {
  try {
    const response = await fetch(`${BACKEND_URL}/status`, {
      cache: "no-store",
    });

    if (!response.ok) {
      return NextResponse.json(
        { status: "unknown", error: "Failed to get status" },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    return NextResponse.json({
      status: "unknown",
      error: `Connection error: ${error}`,
    });
  }
}
