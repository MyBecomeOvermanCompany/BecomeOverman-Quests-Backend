import React from 'react';
import { Grid } from '@mui/material';

function MasonryGrid({ children }) {
  return (
    <Grid container spacing={3}>
      {React.Children.map(children, (child, index) => (
        <Grid item xs={12} sm={6} md={4} lg={3} key={index}>
          {child}
        </Grid>
      ))}
    </Grid>
  );
}

export default MasonryGrid;
