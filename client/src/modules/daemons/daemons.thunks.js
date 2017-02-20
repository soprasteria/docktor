import { withAuth } from '../auth/auth.wrappers';
import { checkHttpStatus, parseJSON, handleError } from '../utils/promises';

import { generateEntitiesThunks } from '../utils/entities';

// Daemons Actions
import DaemonsActions from './daemons.actions';
import DaemonsConstants from './daemons.constants';

/********** Thunk Functions **********/

// Thunk to fetch daemon info
const fetchDaemonInfo = (websocket, daemon, force) => {
  return function (dispatch) {
    dispatch(DaemonsActions.requestDaemonInfo(daemon));
    websocket.send(JSON.stringify({
      'action': DaemonsConstants.REQUEST_DAEMON_INFO,
      'data': {
        'daemon': daemon,
        'force': force
      }
    }));

    // Result is handled by websocket itselve, one time.

  };
};

// Thunk to get all daemons used on a group:
const fetchGroupDaemons = (groupId) => {
  return function (dispatch) {

    dispatch(DaemonsActions.requestAll());

    let request = new Request(`/api/groups/${groupId}/daemons`, withAuth({
      method: 'GET',
    }));
    return fetch(request)
      .then(checkHttpStatus)
      .then(parseJSON)
      .then(response => {
        dispatch(DaemonsActions.receiveSome(response));
      })
      .catch(error => {
        handleError(error, DaemonsActions.invalidRequest, dispatch);
      });
  };
};

export default {
  ...generateEntitiesThunks('daemons'),
  fetchDaemonInfo,
  fetchGroupDaemons
};
