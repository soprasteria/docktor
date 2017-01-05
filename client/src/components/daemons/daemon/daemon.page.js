// React
import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';
import { Scrollbars } from 'react-custom-scrollbars';
import { Form, Input, Button, Dimmer, Loader, Label, Icon, Popup } from 'semantic-ui-react';
import Joi from 'joi-browser';
import UUID from 'uuid-js';

// Thunks / Actions
import SitesThunks from '../../../modules/sites/sites.thunks.js';
import TagsThunks from '../../../modules/tags/tags.thunks.js';
import DaemonsThunks from '../../../modules/daemons/daemons.thunks.js';
import ToastsActions from '../../../modules/toasts/toasts.actions.js';

// Components
import VolumesBox from '../../common/boxes/volumes.box.component.js';
import VariablesBox from '../../common/boxes/variables.box.component.js';
import TagsSelector from '../../tags/tags.selector.component.js';

import { parseError } from '../../../modules/utils/forms.js';

// Style
import './daemon.page.scss';

// Daemon Component
class DaemonComponent extends React.Component {

  state = { errors: { details: [], fields: {} }, daemon: {}, tags:[] }

  schema = Joi.object().keys({
    site: Joi.string().trim().required().label('Site'),
    name: Joi.string().trim().required().label('Name'),
    mountingPoint: Joi.string().trim().required().label('Mounting Point'),
    protocol: Joi.string().trim().required().label('Protocol'),
    host: Joi.string().trim().required().label('Host'),
    port: Joi.number().required().label('Port'),
    timeout: Joi.number().required().label('Timeout')
  })

  componentWillMount = () => {
    this.setState({ daemon: { ...this.props.daemon }, errors: { details: [], fields:{} } });
  }

  componentWillReceiveProps = (nextProps) => {
    this.setState({ daemon: { ...nextProps.daemon }, errors: { details: [], fields:{} } });
  }

  componentDidMount = () => {
    const daemonId = this.props.daemonId;

    // Tags must be fetched before the daemon for the UI to render correctly
    Promise.all([
      this.props.fetchSites(),
      this.props.fetchTags()
    ]).then(() => {
      if (daemonId) {
        // Fetch when known daemon
        this.props.fetchDaemon(daemonId);
      }
    });

    if (!daemonId) {
      // New daemon
      const volumesBox = this.refs.volumes;
      volumesBox.setState({ volumes: [] });
      const variablesBox = this.refs.variables;
      variablesBox.setState({ variables: [] });
      const tagsSelector = this.refs.tags;
      tagsSelector.setState({ tags: [] });
      this.refs.scrollbars && this.refs.scrollbars.scrollTop();
    }
  }

  componentDidUpdate = (prevProps, prevState) => {
    if (prevProps.isFetching) {
      this.refs.scrollbars && this.refs.scrollbars.scrollTop();
    }
  }

  handleChange = (e, { name, value }) => {
    const { daemon, errors } = this.state;
    const state = {
      daemon: { ...daemon, [name]:value },
      errors: { details: [...errors.details], fields: { ...errors.fields } }
    };
    delete state.errors.fields[name];
    this.setState(state);
  }

  isFormValid = () => {
    const { error } = Joi.validate(this.state.daemon, this.schema, { abortEarly: false, allowUnknown: true });
    error && this.setState({ errors: parseError(error) });
    return !Boolean(error);
  }

  onSave(e) {
    e.preventDefault();
    const volumesBox = this.refs.volumes;
    const variablesBox = this.refs.variables;
    const tagsSelector = this.refs.tags;
    // isFormValid validate the form and return the status so all the forms must be validated before doing anything
    let formValid = this.isFormValid() & volumesBox.isFormValid() & variablesBox.isFormValid();
    if (formValid) {
      const daemon = { ...this.state.daemon };
      daemon.volumes = volumesBox.state.volumes;
      daemon.variables = variablesBox.state.variables;
      daemon.tags = tagsSelector.state.tags;
      this.props.onSave(daemon);
    }
  }

  renderSites = (sites, daemon, errors) => {
    const options = sites.map(site => {return { text: site.title, value: site.id };});
    return (
      <Form.Dropdown name='site' label='Site' fluid value={daemon.site} selection placeholder='Select a site...' autoComplete='off' options={options || []}
        required onChange={this.handleChange} loading={!options} error={errors.fields['site']} width='four'
      />
    );
  }

  renderProtocol = (daemon, errors) => {
    const options = [
      { text: 'HTTP', value: 'http' },
      { text: 'HTTPS', value: 'https' }
    ];
    return (
      <Form.Dropdown name='protocol' label='Protocol' fluid value={daemon.protocol} selection placeholder='Select a protocol...' autoComplete='off' options={options}
        required onChange={this.handleChange} error={errors.fields['protocol']} width='three'
      />
    );
  }

  renderCertificates = (daemon) => {
    return (
      <Form.Group widths='three'>
        <Form.TextArea label='CA' name='ca' value={daemon.ca || ''} onChange={this.handleChange}
          rows='10' placeholder='The Certification Authority key Pem file' autoComplete='off'
        />
        <Form.TextArea label='Cert' name='cert' value={daemon.cert || ''} onChange={this.handleChange}
          rows='10' placeholder='The certificate Pem file' autoComplete='off'
        />
        <Form.TextArea label='Key' name='key' value={daemon.key || ''} onChange={this.handleChange}
          rows='10' placeholder='The private key file' autoComplete='off'
        />
      </Form.Group>
    );
  }

  render = () => {
    const { daemon, errors } = this.state;
    const { isFetching, sites, tags } = this.props;
    const certificates = daemon.protocol === 'https';
    const popup = (
      <div>
        Example: <strong>http://host:port/api/v1.x</strong>
        <br/>
        cAdvisor is used to retrieve monitoring stats (CPU, RAM, FS) on host where docker's daemon is running.
        <hr/>
        Docktor recommends to have a cAdvisor instance for each daemon.
      </div>
    );
    return (
      <div className='flex layout vertical start-justified daemon-page'>
        <Scrollbars autoHide ref='scrollbars' className='flex ui dimmable'>
          <div className='flex layout horizontal around-justified'>
            {isFetching && <Dimmer active><Loader size='big' content='Fetching'/></Dimmer>}
            <div className='flex layout vertical start-justified daemon-details'>
              <h1>
                <Link to={'/daemons'}>
                  <Icon name='arrow left' fitted/>
                </Link>
                {this.props.daemon.name || 'New Daemon'}
                <Button size='large' content='Remove' color='red' labelPosition='left' icon='trash'
                  disabled={!daemon.id} onClick={() => this.props.onDelete(daemon)} className='right-floated'
                />
              </h1>
              <Form className='daemon-form'>
                <Input type='hidden' name='created' value={daemon.created || ''} onChange={this.handleChange} />
                <Input type='hidden' name='id' value={daemon.id || ''} onChange={this.handleChange} />

                <Form.Group widths='two'>
                  <Form.Input required label='Name' name='name' value={daemon.name || ''} onChange={this.handleChange}
                    type='text' placeholder='A unique name' autoComplete='off' error={errors.fields['name']}
                  />
                  <Form.TextArea label='Description' name='description' value={daemon.description || ''} onChange={this.handleChange}
                    rows='4' placeholder='A description of the daemon' autoComplete='off'
                  />
                </Form.Group>

                <Form.Group>
                  {this.renderSites(sites, daemon, errors)}
                  <Form.Input required label='Default data mounting point' name='mountingPoint' value={daemon.mountingPoint || ''} onChange={this.handleChange}
                    type='text' placeholder='/data' autoComplete='off' error={errors.fields['mountingPoint']} width='twelve'
                  />
                </Form.Group>

                <Form.Group>
                  <Form.Field width='two'>
                    <Label size='large' className='form-label' content='Tags' />
                  </Form.Field>
                  <Form.Field width='fourteen'>
                    <label>Tags of the daemon</label>
                    <TagsSelector tagsSelectorId={UUID.create(4).hex} selectedTags={daemon.tags || []} tags={tags} ref='tags' />
                  </Form.Field>
                </Form.Group>

                <Form.Group widths='two'>
                  <Form.Field width='two'>
                    <Label size='large' className='form-label' content='cAdvisor' />
                  </Form.Field>
                  <Form.Input label='cAdvisor API URL' name='cadvisorApi' value={daemon.cadvisorApi || ''} onChange={this.handleChange}
                    type='text' autoComplete='off' labelPosition='right corner' width='fourteen'>
                    <input placeholder='http://host:port/api/v1.x' />
                    <Popup trigger={<Label corner='right'><Icon link name='help circle'/></Label>} inverted wide='very'>{popup}</Popup>
                  </Form.Input>
                </Form.Group>

                <Form.Group>
                  <div className='two wide field'>
                    <div className='large ui label form-label'>Docker</div>
                  </div>
                  {this.renderProtocol(daemon, errors)}
                  <Form.Input required label='Hostname' name='host' value={daemon.host || ''} onChange={this.handleChange}
                    type='text' placeholder='Hostname or IP' autoComplete='off' error={errors.fields['host']} width='five'
                  />
                  <Form.Input required label='Port' min='0' name='port' value={daemon.port || ''} onChange={this.handleChange}
                    type='number' placeholder='Port' autoComplete='off' error={errors.fields['port']} width='three'
                  />
                  <Form.Input required label='Timeout' min='0' name='timeout' value={daemon.timeout || ''} onChange={this.handleChange}
                    type='number' placeholder='Timeout' autoComplete='off' error={errors.fields['timeout']} width='three'
                  />
                </Form.Group>
                {certificates && this.renderCertificates(daemon, errors)}
              </Form>

              <VolumesBox volumes={daemon.volumes} ref='volumes'>
                <p>These volumes are used to have common volumes mapping on all services deployed on this daemon. You can add / remove / modify volumes mapping when you deploy a new service on a group.</p>
              </VolumesBox>

              <VariablesBox variables={daemon.variables} ref='variables'>
                <p>These variables are used to have common variables environment into all services deployed on this daemon (Proxy, LDAP,...). You can add / remove / modify variables when you deploy a new service on a group.</p>
              </VariablesBox>

              <div className='flex button-form'>
                <a className='ui fluid button' onClick={event => this.onSave(event)}>Save</a>
              </div>
            </div>
          </div>
        </Scrollbars>
      </div>
    );
  }
}
DaemonComponent.propTypes = {
  daemon: React.PropTypes.object,
  isFetching: React.PropTypes.bool,
  daemonId: React.PropTypes.string,
  sites: React.PropTypes.array,
  tags: React.PropTypes.object,
  fetchDaemon: React.PropTypes.func.isRequired,
  fetchSites: React.PropTypes.func.isRequired,
  fetchTags: React.PropTypes.func.isRequired,
  onSave: React.PropTypes.func,
  onDelete: React.PropTypes.func
};

// Function to map state to container props
const mapStateToProps = (state, ownProps) => {
  const paramId = ownProps.params.id;
  const daemons = state.daemons;
  const daemon = daemons.selected;
  const emptyDaemon = { volumes: [], variables: [], tags: [] };
  const isFetching = paramId && (paramId !== daemon.id || (daemon.id ? daemon.isFetching : true));
  const sites = Object.values(state.sites.items);
  return {
    daemon: daemons.items[paramId] || emptyDaemon,
    isFetching,
    daemonId: paramId,
    sites,
    tags: state.tags
  };
};

// Function to map dispatch to container props
const mapDispatchToProps = (dispatch) => {
  return {
    fetchDaemon: id => dispatch(DaemonsThunks.fetchDaemon(id)),
    fetchSites: () => dispatch(SitesThunks.fetchIfNeeded()),
    fetchTags: () => dispatch(TagsThunks.fetchIfNeeded()),
    onSave: daemon => dispatch(DaemonsThunks.saveDaemon(daemon)),
    onDelete: daemon => {
      const callback = () => dispatch(DaemonsThunks.deleteDaemon(daemon.id));
      dispatch(ToastsActions.confirmDeletion(daemon.name, callback));
    }
  };
};

// Redux container to Sites component
const DaemonPage = connect(
  mapStateToProps,
  mapDispatchToProps
)(DaemonComponent);

export default DaemonPage;
