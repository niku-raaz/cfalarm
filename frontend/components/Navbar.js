import React from 'react';
import Link from 'next/link'; // Use Next.js Link for client-side navigation

export default function Navbar() {
  return (
    <nav className="navbar">
      <h1>CFAlarm</h1>
      <div className="links">
        <Link href="/">Home</Link>
        <Link href="/contests">Contests</Link>
        <Link href="/profile">Profile</Link> {/* Add Profile Link */}
        <Link href="/login">Login</Link>
      </div>
    </nav>
  );
}