// React
import React from 'react';
import { Map, TileLayer, Marker, Popup } from 'react-leaflet';
import { connect } from 'react-redux';

// Actions for redux container
import SitesThunks from './sites.thunks.js';
import ToastsActions from '../toasts/toasts.actions.js';
import ModalActions from '../modal/modal.actions.js';

// Style
import 'leaflet/dist/leaflet.css';
import './sites.component.scss';

//Site Component using react-leaflet
class SitesComponent extends React.Component {
  constructor() {
    super();
    this.initPosition = { lat: 45, lng: 5, zoom: 4 };
  }

  componentDidUpdate() {
    setTimeout(() => {
      this.refs.map1.leafletElement.invalidateSize(false);
    }, 300); // Adjust timeout to tab transition
  }

  openModalNewSite(onCreate) {
    return e => {
      onCreate(e.latlng);
    };
  }

  openModalEditSite(onEdit, site) {
    return () => {
      onEdit(site);
    };
  }

  render() {
    const initPosition = [this.initPosition.lat, this.initPosition.lng];
    const sites = Object.values(this.props.sites.items);
    const fetching = this.props.sites.isFetching;
    const onDelete = this.props.onDelete;
    const onCreate = this.props.onCreate;
    const onEdit = this.props.onEdit;
    return (
      <div className='flex-2 self-stretch map-container layout horizontal center-center'>
        <Map ref='map1' className='flex self-stretch map' center={initPosition} zoom={this.initPosition.zoom} onClick={this.openModalNewSite(onCreate)}>
          {(fetching => {
            if (fetching) {
              return (
                <div className='ui active inverted dimmer'>
                  <div className='ui text loader'>Fetching</div>
                </div>
              );
            }
          })(fetching)}
          <TileLayer
            attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
            url='http://{s}.tile.osm.org/{z}/{x}/{y}.png'
            />
          {sites.map(site => {
            const sitePosition = [site.latitude, site.longitude];
            return (
              <Marker key={site.id} position={sitePosition}>
                <Popup>
                  <div>{site.title}{' '}
                    <i onClick={this.openModalEditSite(onEdit, site)} className='blue write link icon'></i>
                    <i onClick={() => onDelete(site)} className='red trash link icon'></i>
                  </div>
                </Popup>
              </Marker>
            );
          })}
        </Map>
      </div>
    );
  }
}
SitesComponent.propTypes = {
  sites: React.PropTypes.object,
  onCreate: React.PropTypes.func,
  onEdit: React.PropTypes.func,
  onDelete: React.PropTypes.func
};


// Function to map state to container props
const mapStateToSitesProps = (state) => {
  return { sites: state.sites };
};

// Function to map dispatch to container props
const mapDispatchToSitesProps = (dispatch) => {
  return {
    onDelete: site => {
      const callback = () => dispatch(SitesThunks.deleteSite(site.id));
      dispatch(ToastsActions.confirmDeletion(site.title, callback));
    },
    onCreate: position => {
      const callback = (siteForm) => dispatch(SitesThunks.saveSite(siteForm));
      dispatch(ModalActions.openNewSiteModal(position, callback));
    },
    onEdit: site => {
      const callback = (siteForm) => dispatch(SitesThunks.saveSite(siteForm));
      dispatch(ModalActions.openEditSiteModal(site, callback));
    }
  };
};

// Redux container to Sites component
const Sites = connect(
  mapStateToSitesProps,
  mapDispatchToSitesProps
)(SitesComponent);

export default Sites;
