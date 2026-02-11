import { NextRequest, NextResponse } from "next/server";

export const dynamic = "force-dynamic";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function GET(request: NextRequest) {
  try {
    const path = request.nextUrl.searchParams.get("path");
    if (!path) {
      return NextResponse.json(
        { error: "path parameter required" },
        { status: 400 }
      );
    }

    const response = await fetch(
      `${BACKEND_URL}/check-dir?path=${encodeURIComponent(path)}`,
      {
        cache: "no-store",
      }
    );

    if (!response.ok) {
      return NextResponse.json(
        { error: "Failed to check directory" },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    // If backend is unreachable, assume dir is ok (we can't check)
    return NextResponse.json({ exists: false, empty: true });
  }
}
