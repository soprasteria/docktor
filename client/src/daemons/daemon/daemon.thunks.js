import 'babel-polyfill';
import fetch from 'isomorphic-fetch';
import { withAuth } from '../../auth/auth.wrappers.js';
import { checkHttpStatus, parseJSON, handleError } from '../../utils/utils.js';

// Daemon Actions
import DaemonActions from './daemon.actions.js';

/********** Thunk Functions **********/

// Thunk to fetch daemons
const fetchDaemon = (id) => {
  return function (dispatch) {

    dispatch(DaemonActions.requestDaemon(id));

    return fetch(`/api/daemons/${id}`, withAuth({ method:'GET' }))
      .then(checkHttpStatus)
      .then(parseJSON)
      .then(response => {
        dispatch(DaemonActions.receiveDaemon(response));
      })
      .catch(error => {
        handleError(error, DaemonActions.invalidRequestDaemon, dispatch);
      });
  };
};

// Thunk to save daemons
const saveDaemon = (form) => {

  let daemon = Object.assign({}, form);
  daemon.port = parseInt(daemon.port);
  daemon.timeout = parseInt(daemon.timeout);
  daemon.created = daemon.created ? daemon.created : new Date();
  const id = form.id ? form.id : -1;
  return function (dispatch) {

    dispatch(DaemonActions.requestSaveDaemon(daemon));

    let request = new Request('/api/daemons/' + id, withAuth({
      method: 'PUT',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(daemon)
    }));

    return fetch(request)
      .then(checkHttpStatus)
      .then(parseJSON)
      .then(response => {
        dispatch(DaemonActions.savedDaemon(response));
      })
      .catch(error => {
        handleError(error, DaemonActions.invalidRequestDaemon, dispatch);
      });
  };
};

 export default {
   fetchDaemon,
   saveDaemon
 };
