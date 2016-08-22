// Imports for fetch API
import 'babel-polyfill';

// Daemon Actions
import {
  requestLogin,
  loginError,
  receiveLogin,
  requestLogout,
  receiveLogout
} from './auth.actions.js';

// Calls the API to get a token and
// dispatches actions along the way
export function loginUser(creds) {

  let config = {
    method: 'POST',
    headers: { 'Content-Type':'application/x-www-form-urlencoded' },
    body: `username=${creds.username}&password=${creds.password}`
  };

  return dispatch => {
    // We dispatch requestLogin to kickoff the call to the API
    dispatch(requestLogin());

    return fetch('/create-token', config)
      .then(response =>
        response.json().then(
          user => ({ user, response })
        )
      ).then(({ user, response }) =>  {
        if (!response.ok) {
          // If there was a problem, we want to
          // dispatch the error condition
          dispatch(loginError(user.message));
          return Promise.reject(user);
        } else {
          // If login was successful, set the token in local storage
          localStorage.setItem('id_token', user.id_token);
          // Dispatch the success action
          dispatch(receiveLogin(user));
        }
      }).catch(err => console.log('Error: ', err));
  };
}

// Logs the user out
export function logoutUser() {
  return dispatch => {
    dispatch(requestLogout());
    localStorage.removeItem('id_token');
    dispatch(receiveLogout());
  };
}
