import Navbar from '../components/Navbar';
import { useRouter } from 'next/router'; // Import the router hook

export default function Login() {
  const router = useRouter(); // Initialize the router

  const handleLogin = () => {
    // This line is correct. It sends the user to the backend to start the Google login flow.
    window.location.href = 'http://localhost:8080/api/auth/google';
    
  };

  return (
    <div>
      <Navbar />
      <div style={{ padding: '2rem' }}>
        <h2>Login via Google</h2>
        <button onClick={handleLogin} style={{ padding: '0.5rem 1rem', fontSize: '1rem' }}>
          Login with Google
        </button>
      </div>
    </div>
  );
}