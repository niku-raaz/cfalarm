import Navbar from '../components/Navbar';

export default function Login() {
  const handleLogin = () => {
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
