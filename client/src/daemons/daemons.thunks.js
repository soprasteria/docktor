// Imports for fetch API
import 'babel-polyfill';
import fetch from 'isomorphic-fetch';
import { withAuth } from '../auth/auth.wrappers.js';
import { checkHttpStatus, parseJSON, dispatchError } from '../utils/utils.js';

// Daemon Actions
import {
  requestAllDaemons,
  receiveDaemons,
  invalidRequestDaemons
} from './daemons.actions.js';

/********** Thunk Functions **********/

// Thunk to fetch daemons
export function fetchDaemons() {
  return function (dispatch) {

    dispatch(requestAllDaemons());

    return fetch('/api/daemons', withAuth({ method:'GET' }))
      .then(checkHttpStatus)
      .then(parseJSON)
      .then(response => {
        dispatch(receiveDaemons(response));
      })
      .catch(error => {
        handleError(error, invalidRequestSites, dispatch);
      });
  };
}


/********** Helper Functions **********/

// Check that if daemons should be fetched
function shouldFetchDaemons(state) {
  const daemons = state.daemons;
  if (!daemons || daemons.didInvalidate) {
    return true;
  } else if (daemons.isFetching) {
    return false;
  } else {
    return true;
  }
}

// Thunk to fech daemons only if needed
export function fetchDaemonsIfNeeded() {

  return (dispatch, getState) => {
    if (shouldFetchDaemons(getState())) {
      return dispatch(fetchDaemons());
    } else {
      return Promise.resolve();
    }
  };
}