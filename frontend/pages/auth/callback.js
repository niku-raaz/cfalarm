import { useEffect } from 'react';
import { useRouter } from 'next/router';

export default function AuthCallback() {
  const router = useRouter();

  useEffect(() => {
    // This effect runs when the page loads.
    // We check if the 'token' is present in the URL query parameters.
    if (router.query.token) {
      const { token } = router.query;
      
      // Store the token in localStorage for future authenticated requests.
      localStorage.setItem('authToken', token);
      
      // Redirect the user to their profile page.
      // The `replace` method is used so the user can't click "back" to this page.
      router.replace('/profile');
    }
  }, [router.query, router]);

  // You can show a loading message while the redirect happens.
  return (
    <div style={{ padding: '2rem' }}>
      <p>Please wait, we are logging you in...</p>
    </div>
  );
}