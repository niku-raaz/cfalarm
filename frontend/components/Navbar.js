import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';

export default function Navbar() {
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('authToken');
    if (token) {
      setIsLoggedIn(true);
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    setIsLoggedIn(false);
    router.push('/login');
  };

  return (
    <nav className="navbar">
      <h1>CFAlarm</h1>
      <div className="links">
        <Link href="/">Home</Link>
        
        {/* --- UPDATED LOGIC --- */}
        {isLoggedIn ? (
          <>
            <Link href="/practice">Practice</Link> {/* New Practice Link */}
            <Link href="/profile">Profile</Link>
            <button onClick={handleLogout} className="logout-button" style={{background: 'none', border: 'none', color: 'white', fontWeight: 'bold', cursor: 'pointer', fontSize: '1rem' }}>
              Logout
            </button>
          </>
        ) : (
          <Link href="/login">Login</Link>
        )}
        {/* --- END OF UPDATED LOGIC --- */}
      </div>
    </nav>
  );
}