import { generateEntitiesConstants } from '../utils/entities.js';

export default {
  ...generateEntitiesConstants('tags'),
  REQUEST_CREATE_TAG: 'REQUEST_CREATE_TAG',
  RECEIVE_TAG_CREATED: 'RECEIVE_TAG_CREATED',
  CREATE_TAG_INVALID: 'CREATE_TAG_INVALID',
  REQUEST_DELETE_TAG: 'REQUEST_DELETE_TAG',
  RECEIVE_TAG_DELETED: 'RECEIVE_TAG_DELETED',
  DELETE_TAG_INVALID: 'DELETE_TAG_INVALID',
  REQUEST_SAVE_TAG: 'REQUEST_SAVE_TAG',
  RECEIVE_TAG_SAVED: 'RECEIVE_TAG_SAVED',
  SAVE_TAG_INVALID: 'SAVE_TAG_INVALID'
};