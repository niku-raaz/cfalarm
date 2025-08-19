import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Navbar from '../components/Navbar';
import axios from 'axios';

export default function Profile() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [codeforcesId, setCodeforcesId] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    const fetchProfile = async () => {
      // 1. Log the token right after getting it from storage
      const token = localStorage.getItem('authToken');
      console.log('Token from localStorage:', token);

      if (!token) {
        console.log('No token found, redirecting to login.');
        router.replace('/login');
        return;
      }

      try {
        // 2. Log the headers just before making the API call
        const headers = { Authorization: `Bearer ${token}` };
        console.log('Sending request with headers:', headers);

        const response = await axios.get('http://localhost:8080/api/user/me', { headers });
        
        const userData = response.data;
        setUser(userData);
        setCodeforcesId(userData.CodeforcesID || '');
        
        if (!userData.CodeforcesID) {
          setIsEditing(true);
        }

      } catch (err) {
        // 3. Log the error if the API call fails
        console.error('Error fetching profile:', err);
        setError('Failed to fetch profile. Please log in again.');
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
        'http://localhost:8080/api/user/verify-cf-handle',
        { CodeforcesID: codeforcesId },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setMessage(response.data.message);
      setIsEditing(false);
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Failed to save handle.';
      setError(errorMessage);
    } finally {
      setIsSaving(false);
    }
  };
  
  const handleEditClick = () => {
    setIsEditing(true);
    setMessage('');
  };

  if (loading) {
    return (
      <div>
        <Navbar />
        <div style={{ padding: '2rem' }}>Loading...</div>
      </div>
    );
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

            <div style={{ marginTop: '2rem' }}>
              <p><strong>Codeforces Handle:</strong></p>
              {isEditing ? (
                <form onSubmit={handleVerifyAndSave} style={{ display: 'flex', alignItems: 'center' }}>
                  <input 
                    type="text" 
                    value={codeforcesId} 
                    onChange={(e) => setCodeforcesId(e.target.value)} 
                    placeholder="Your Codeforces Handle" 
                    style={{ padding: '0.5rem', width: '300px' }} 
                  />
                  <button 
                    type="submit" 
                    style={{ padding: '0.5rem 1rem', marginLeft: '1rem' }} 
                    disabled={isSaving}
                  >
                    {isSaving ? 'Verifying...' : 'Verify & Save'}
                  </button>
                </form>
              ) : (
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <p style={{ margin: 0, fontSize: '1.1rem' }}>{codeforcesId}</p>
                  <button 
                    onClick={handleEditClick}
                    style={{ padding: '0.5rem 1rem', marginLeft: '1rem' }} 
                  >
                    Update
                  </button>
                </div>
              )}
            </div>
            
            {message && <p style={{ color: 'green', marginTop: '1rem' }}>{message}</p>}
            {error && <p style={{ color: 'red', marginTop: '1rem' }}>{error}</p>}
          </div>
        )}
      </div>
    </div>
  );
}