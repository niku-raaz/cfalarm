import React from 'react';


export default function Navbar() {
  return (
    <nav className="navbar">
      <h1>CFAlarm</h1>
      <div className="links">
        <a href="/">Home</a>
        <a href="/contests">Contests</a>
        <a href="/login">Login</a>
      </div>
    </nav>
  );
}
