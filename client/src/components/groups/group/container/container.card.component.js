// React
import React from 'react';
import { Link } from 'react-router';

// Style
import './container.card.component.scss';

// ContainerCard Component
class ContainerCard extends React.Component {
  render() {
    const { container, daemons } = this.props;
    const isFetching = false;
    const daemon = daemons.find(daemon => container.daemonId === daemon.value) || { 'name': 'Unknown' };
    const statusMessage = `Container is up on daemon '${daemon.name}'`;
    return (
      <div className='container'>
        <div className='ui card'>
          <div className='content'>
            <div className='header'>
              <Link className='header' to={'containers/' + container.id} title={container.serviceTitle}>
                {container.serviceTitle}
              </Link>
            </div>
            <div title={statusMessage} className={'ui top right attached label green'}>
              <i className='refresh icon' />UP
            </div>
            <div className='meta'>{container.name}</div>
            <div className='description'>{container.image}</div>
          </div>
          <div className='ui bottom attached buttons'>
            <div className='ui icon button'><i className='stop icon' />Stop</div>
            <div className='ui icon button'><i className='play icon' />Start</div>
            <div className='ui icon button'><i className='repeat icon' />Restart</div>
            <div className='ui icon button'><i className='cloud upload icon' />Deploy</div>
          </div>
        </div>
      </div>
    );
  }
}
ContainerCard.propTypes = { container: React.PropTypes.object, daemons: React.PropTypes.array };

export default ContainerCard;
