import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Navbar from '../components/Navbar';
import axios from 'axios';

export default function Profile() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [codeforcesId, setCodeforcesId] = useState('');
  const [apiKey, setApiKey] = useState('');
  const [apiSecret, setApiSecret] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [isVerified, setIsVerified] = useState(false);

  useEffect(() => {
    const fetchProfile = async () => {
      const token = localStorage.getItem('authToken');
      if (!token) {
        router.replace('/login');
        return;
      }
      try {
        const response = await axios.get('http://localhost:8080/api/user/me', {
          headers: { Authorization: `Bearer ${token}` },
        });
        const userData = response.data;
        setUser(userData);
        setCodeforcesId(userData.CodeforcesID || '');
        // We still fetch these to check if they exist, but they won't be shown initially
        if (userData.CfApiKey) {
          setIsVerified(true);
        }
      } catch (err) {
        setError('Failed to fetch profile. Please log in again.');
        localStorage.removeItem('authToken');
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, [router]);

  const handleVerifyAndSave = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('authToken');
    if (!token) return;

    setIsSaving(true);
    setMessage('');
    setError('');

    try {
      const response = await axios.post(
        'http://localhost:8080/api/user/verify-cf-keys',
        {
          codeforcesId: codeforcesId,
          apiKey: apiKey,
          apiSecret: apiSecret,
        },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setMessage(response.data.message);
      setIsVerified(true); 
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'An error occurred during verification.';
      setError(errorMessage);
    } finally {
      setIsSaving(false);
    }
  };

  // --- NEW HANDLER FUNCTION ---
  // This function will clear the API keys and show the form.
  const handleUpdateClick = () => {
    setApiKey('');
    setApiSecret('');
    setIsVerified(false);
  };

  if (loading) {
    return ( <div><Navbar /><div style={{ padding: '2rem' }}>Loading...</div></div> );
  }

  return (
    <div>
      <Navbar />
      <div style={{ padding: '2rem' }}>
        {user && (
          <div>
            <h2>User Profile</h2>
            <p><strong>Name:</strong> {user.Name}</p>
            <p><strong>Email:</strong> {user.Email}</p>

            {isVerified ? (
              <div style={{ marginTop: '2rem' }}>
                <p><strong>Codeforces Handle:</strong> {codeforcesId}</p>
                <p style={{ color: 'green' }}>✓ Credentials Verified</p>
                {/* Use the new handler function here */}
                <button
                  onClick={handleUpdateClick}
                  style={{ padding: '0.5rem 1rem', fontSize: '1rem', marginTop: '1rem' }}
                >
                  Update Credentials
                </button>
              </div>
            ) : (
              <>
                <p style={{ marginTop: '2rem', fontStyle: 'italic' }}>
                  Generate your API keys at: 
                  <a href="https://codeforces.com/settings/api" target="_blank" rel="noopener noreferrer">
                    codeforces.com/settings/api
                  </a>
                </p>

                <form onSubmit={handleVerifyAndSave}>
                  <div style={{ marginTop: '1rem' }}>
                    <label><strong>Codeforces Handle:</strong></label><br />
                    <input type="text" value={codeforcesId} onChange={(e) => setCodeforcesId(e.target.value)} placeholder="Your Codeforces Handle" style={{ padding: '0.5rem', marginTop: '0.5rem', width: '300px' }} />
                  </div>
                  <div style={{ marginTop: '1rem' }}>
                    <label><strong>API Key:</strong></label><br />
                    <input type="text" value={apiKey} onChange={(e) => setApiKey(e.target.value)} placeholder="Your Codeforces API Key" style={{ padding: '0.5rem', marginTop: '0.5rem', width: '300px' }} />
                  </div>
                  <div style={{ marginTop: '1rem' }}>
                    <label><strong>API Secret:</strong></label><br />
                    <input type="password" value={apiSecret} onChange={(e) => setApiSecret(e.target.value)} placeholder="Your Codeforces API Secret" style={{ padding: '0.5rem', marginTop: '0.5rem', width: '300px' }} />
                  </div>
                  <button type="submit" style={{ padding: '0.5rem 1rem', fontSize: '1rem', marginTop: '1rem' }} disabled={isSaving}>
                    {isSaving ? 'Verifying...' : 'Verify & Save Credentials'}
                  </button>
                </form>
              </>
            )}
            
            {message && <p style={{ color: 'green', marginTop: '1rem' }}>{message}</p>}
            {error && <p style={{ color: 'red', marginTop: '1rem' }}>{error}</p>}
          </div>
        )}
      </div>
    </div>
  );
}