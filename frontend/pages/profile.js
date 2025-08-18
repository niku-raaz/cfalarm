import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Navbar from '../components/Navbar';
import axios from 'axios';

export default function Profile() {
  const router = useRouter();
  // State to hold user data
  const [user, setUser] = useState(null);
  // State to manage the Codeforces ID input field
  const [codeforcesId, setCodeforcesId] = useState('');
  // State for loading and error messages
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    // This function fetches the user's profile data from the backend.
    const fetchProfile = async () => {
      // Retrieve the auth token from localStorage.
      const token = localStorage.getItem('authToken');

      // If there's no token, the user is not logged in. Redirect them.
      if (!token) {
        router.replace('/login');
        return;
      }

      try {
        // Make an authenticated GET request to the /api/user/me endpoint.
        const response = await axios.get('http://localhost:8080/api/user/me', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        
        // On success, update the component's state with the user data.
        setUser(response.data);
        setCodeforcesId(response.data.CodeforcesID || ''); // Pre-fill the input
      } catch (err) {
        // If the request fails (e.g., token is invalid), set an error message.
        setError('Failed to fetch profile. Please log in again.');
        localStorage.removeItem('authToken'); // Clear the bad token
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [router]);

  // This function is called when the user saves their Codeforces ID.
  const handleUpdateProfile = async (e) => {
    e.preventDefault(); // Prevent the default form submission behavior
    const token = localStorage.getItem('authToken');
    if (!token) return;

    try {
      // Make an authenticated PUT request to update the user's profile.
      await axios.put(
        'http://localhost:8080/api/user/me',
        {
          CodeforcesID: codeforcesId,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      // On success, show a confirmation message.
      setMessage('Profile updated successfully!');
    } catch (err) {
      setError('Failed to update profile.');
    }
  };
  
  // Display a loading message while data is being fetched.
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
        {error && <p style={{ color: 'red' }}>{error}</p>}
        {user && (
          <div>
            <h2>User Profile</h2>
            <p><strong>Name:</strong> {user.Name}</p>
            <p><strong>Email:</strong> {user.Email}</p>

            <form onSubmit={handleUpdateProfile}>
              <div style={{ marginTop: '1rem' }}>
                <label htmlFor="codeforcesId">
                  <strong>Codeforces ID:</strong>
                </label>
                <br />
                <input
                  id="codeforcesId"
                  type="text"
                  value={codeforcesId}
                  onChange={(e) => setCodeforcesId(e.target.value)}
                  placeholder="Enter your Codeforces handle"
                  style={{ padding: '0.5rem', marginTop: '0.5rem', width: '250px' }}
                />
              </div>
              <button
                type="submit"
                style={{ padding: '0.5rem 1rem', fontSize: '1rem', marginTop: '1rem' }}
              >
                Save
              </button>
            </form>
            {message && <p style={{ color: 'green', marginTop: '1rem' }}>{message}</p>}
          </div>
        )}
      </div>
    </div>
  );
}