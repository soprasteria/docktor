// React
import React from 'react';
import PropTypes from 'prop-types';

import Box from './box/box.component';

// VolumesBox is a list of docker volumes
class VolumesBox extends React.Component {

  state = { volumes: [] }

  componentWillMount = () => {
    this.setState({ volumes: this.props.volumes });
  }

  componentWillReceiveProps = (nextProps) => {
    this.setState({ volumes: nextProps.volumes });
  }

  isFormValid = () => {
    return this.refs.volumesBox.isFormValid();
  }

  onChangeVolumes = (volumes) => {
    this.state.volumes = volumes;
  }

  render = () => {
    const form = { fields:[] };
    const allowEmpty = this.props.allowEmpty;

    form.getTitle = (volume) => {
      const external = volume.external || (allowEmpty ? '<Default Value>' : '' );
      const rights = volume.rights || 'rw';
      return '-v ' + external + ':' + volume.internal + ':' + rights;
    };

    form.fields.push({
      name: 'external',
      label: allowEmpty ? 'Default Value' : 'External Volume',
      placeholder: 'The default volume on host',
      class: 'five wide',
      required: !allowEmpty
    });

    form.fields.push({
      name: 'internal',
      label: 'Internal Volume',
      placeholder: 'The volume inside the container',
      class: 'five wide',
      required: true
    });
    form.fields.push({
      name: 'rights',
      label: 'Rights',
      placeholder: 'Select rights',
      class: 'three wide',
      required: true,
      options: [
        { value:'ro', name:'Read-only' },
        { value:'rw', name:'Read-write' }
      ],
      default: 'rw',
      type: 'select'
    });

    form.fields.push({
      name: 'description',
      label: 'Description',
      placeholder: 'Describe this volume',
      class: 'three wide',
      type: 'textarea',
      rows: 2
    });

    return (
      <Box
        ref='volumesBox'
        icon='folder open'
        title='Volumes' form={form}
        lines={this.props.volumes}
        stacked={this.props.stacked}
        onChange={this.onChangeVolumes}>
        {this.props.children || ''}
      </Box>
    );
  }
};

VolumesBox.propTypes = {
  volumes: PropTypes.array,
  allowEmpty: PropTypes.bool,
  stacked: PropTypes.bool,
  children: PropTypes.oneOfType([
    PropTypes.array,
    PropTypes.element
  ])
};

export default VolumesBox;
