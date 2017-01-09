// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { Scrollbars } from 'react-custom-scrollbars';
import UUID from 'uuid-js';

// Thunks / Actions
import TagsThunks from '../../../modules/tags/tags.thunks.js';
import DaemonsThunks from '../../../modules/daemons/daemons.thunks.js';
import UsersThunks from '../../../modules/users/users.thunks.js';

import GroupsThunks from '../../../modules/groups/groups.thunks.js';
import ToastsActions from '../../../modules/toasts/toasts.actions.js';

// Components
import FilesystemsBox from '../../common/boxes/filesystems.box.component.js';
import MembersBox from '../../common/boxes/members.box.component.js';
import TagsSelector from '../../tags/tags.selector.component.js';

// Selectors
import { getDaemonsAsFSOptions } from '../../../modules/daemons/daemons.selectors.js';
import { getUsersAsOptions } from '../../../modules/users/users.selectors.js';

// Style
import './group.edit.page.scss';

// Group Component for edition
class GroupEditComponent extends React.Component {

  constructor(props) {
    super(props);
    this.state = { ...props.group };
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ ...nextProps.group });
  }

  componentDidMount() {
    const groupId = this.props.groupId;

    // Tags must be fetched before the group for the UI to render correctly
    Promise.all([
      this.props.fetchTags(),
      this.props.fetchDaemons(),
      this.props.fetchUsers()
    ]).then(() => {
      if (groupId) {
        // Fetch when known group
        this.props.fetchGroup(groupId);
      }
    });

    if (!groupId) {
      // New group
      $('.ui.form.group-form').form('clear');
      const tagsSelector = this.refs.tags;
      tagsSelector.setState({ tags: [] });
      this.refs.scrollbars.scrollTop();
    }
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevProps.isFetching) {
      this.refs.scrollbars.scrollTop();
    }
  }

  onChangeProperty(value, property) {
    this.setState({ [property]: value });
  }

  isFormValid() {
    const settings = {
      fields: {
        title: 'empty',
        description: 'empty',
        portminrange: 'empty',
        portmaxrange: 'empty'
      }
    };
    $('.ui.form.group-form').form(settings);
    $('.ui.form.group-form').form('validate form');
    return $('.ui.form.group-form').form('is valid');
  }

  onSave(event) {
    event.preventDefault();
    const tagsSelector = this.refs.tags;
    const filesystemsBox = this.refs.filesystemsBox;
    const membersBox = this.refs.membersBox;
    // isFormValid validate the form and return the status so all the forms must be validated before doing anything

    let formValid = filesystemsBox.isFormValid() & membersBox.isFormValid() & this.isFormValid();
    if (formValid) {
      const group = { ...this.state };
      group.tags = tagsSelector.state.tags;
      group.filesystems = filesystemsBox.state.filesystems;
      group.members = membersBox.state.members;
      this.props.onSave(group);
    }
  }

  render() {
    const group = this.state;
    const isFetching = this.props.isFetching;
    const daemons = this.props.daemons;
    const tags = this.props.tags;
    const users = this.props.users;
    return (
      <div className='flex layout vertical start-justified group-page'>
        <Scrollbars ref='scrollbars' className='flex ui dimmable'>
          <div className='flex layout horizontal around-justified'>
            {
              isFetching ?
                <div className='ui active dimmer'>
                  <div className='ui text loader'>Fetching</div>
                </div>
                :
                <div className='flex layout vertical start-justified group-details'>
                  <h1>
                    <Link to={group.id ? `/groups/${group.id}` : '/groups'}>
                      <i className='arrow left icon'/>
                    </Link>
                    {this.props.group.title || 'New Group'}
                    <button disabled={!group.id} onClick={() => this.props.onDelete(group)} className='ui red labeled icon button right-floated'>
                      <i className='trash icon'/>Remove
                    </button>
                  </h1>
                  <form className='ui form group-form'>
                    <input type='hidden' name='created' value={group.created || ''} onChange={event => this.onChangeProperty(event.target.value, 'created')} />
                    <input type='hidden' name='id' value={group.id || ''} onChange={event => this.onChangeProperty(event.target.value, 'id')} />
                    <div className='field required'>
                      <label>Title</label>
                      <input type='text' name='title' value={group.title || ''} onChange={event => this.onChangeProperty(event.target.value, 'title')}
                        placeholder='A unique name' autoComplete='off' />
                    </div>
                    <div className='field'>
                      <label>Description</label>
                      <textarea rows='4' name='description' value={group.description || ''} onChange={event => this.onChangeProperty(event.target.value, 'description')}
                        placeholder='A description of the group' autoComplete='off' />
                    </div>
                    <div className='fields'>
                      <div className='two wide field'>
                        <div className='large ui label form-label'>Tags</div>
                      </div>
                      <div className='fourteen wide field'>
                        <label>Tags of the group</label>
                        <TagsSelector tagsSelectorId={UUID.create(4).hex} selectedTags={group.tags || []} tags={tags} ref='tags' />
                      </div>
                    </div>
                  </form>
                  <FilesystemsBox filesystems={group.filesystems} daemons={daemons} ref='filesystemsBox' boxId={UUID.create(4).hex}>
                    <p>Monitoring filesystem is only available if selected daemon has cAdvisor deployed on it and configured on Docktor.</p>
                  </FilesystemsBox>
                  <MembersBox members={group.members} users={users} ref='membersBox' boxId={UUID.create(4).hex}>
                    <p>Members of groups are able to see it and interact with containers.</p>
                    <ul>
                      <li>Moderators are able to add other members and to interact with services (stop/start)</li>
                      <li>Simple members are only able to see the group and instanciated services</li>
                    </ul>
                  </MembersBox>
                  <div className='flex button-form'>
                    <a className='ui fluid button' onClick={event => this.onSave(event)}>Save</a>
                  </div>
                </div>
            }
          </div>
        </Scrollbars>
      </div>
    );
  }
}
GroupEditComponent.propTypes = {
  group: React.PropTypes.object,
  isFetching: React.PropTypes.bool,
  groupId: React.PropTypes.string,
  daemons: React.PropTypes.array,
  users: React.PropTypes.array,
  tags: React.PropTypes.object,
  fetchGroup: React.PropTypes.func.isRequired,
  fetchDaemons: React.PropTypes.func.isRequired,
  fetchTags: React.PropTypes.func.isRequired,
  fetchUsers: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func,
  onDelete: React.PropTypes.func
};

// Function to map state to container props
const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const groups = state.groups;
  const group = groups.selected;
  const emptyGroup = { tags: [], filesystems: [], members: [] };
  const daemons = getDaemonsAsFSOptions(state.daemons.items) || [];
  const users = getUsersAsOptions(state.users.items) || [];
  const isFetching = paramId && (paramId !== group.id || (group.id ? group.isFetching : true));
  return {
    group: groups.items[paramId] || emptyGroup,
    isFetching,
    groupId: paramId,
    tags: state.tags,
    daemons,
    users
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    fetchGroup: (id) => dispatch(GroupsThunks.fetchGroup(id)),
    fetchDaemons: () => dispatch(DaemonsThunks.fetchIfNeeded()),
    fetchUsers: () => dispatch(UsersThunks.fetchIfNeeded()),
    fetchTags: () => dispatch(TagsThunks.fetchIfNeeded()),
    onSave: (group) => dispatch(GroupsThunks.saveGroup(group)),
    onDelete: group => {
      const callback = () => dispatch(GroupsThunks.deleteGroup(group.id));
      dispatch(ToastsActions.confirmDeletion(group.title, callback));
    }
  };
};

// Redux container to Sites component
const GroupEditPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(GroupEditComponent);

export default GroupEditPage;