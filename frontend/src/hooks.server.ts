export const handle = async ({ event, resolve }) => {
  const token = event.cookies.get('access_token');

  if (token) {
    try {
      const verifyRes = await fetch('http://localhost:8080/auth/verify', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
    

      if (verifyRes.ok) {
        const data = await verifyRes.json();
        event.locals.user = {
          ...data,token
        };
	
      }
    } catch (err) {
      console.error('Token verification failed:', err);
    }
  }

  return resolve(event);
};
