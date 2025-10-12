import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Box,
  AppBar,
  Toolbar,
  IconButton,
  Grid,
  Card,
  CardContent,
  Alert,
  CircularProgress,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  LinearProgress,
} from '@mui/material';
import {
  Logout as LogoutIcon,
  Dashboard as DashboardIcon,
  Memory as MemoryIcon,
  Storage as StorageIcon,
  Computer as ComputerIcon,
  AccountCircle as AccountCircleIcon,
  AdminPanelSettings as AdminIcon,
  Info as InfoIcon,
  Image as ImageIcon,
  Folder as FolderIcon,
  ViewInAr as ViewInArIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import { useAuth } from '../context/AuthContext';
import { systemApi } from '../services/api';

const SystemInfoPage = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [systemInfo, setSystemInfo] = useState(null);
  const [diskUsage, setDiskUsage] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    setError('');
    try {
      const [sysInfo, diskInfo] = await Promise.all([
        systemApi.getSystemInfo(),
        systemApi.getDiskUsage(),
      ]);
      setSystemInfo(sysInfo);
      setDiskUsage(diskInfo);
    } catch (err) {
      setError('Failed to load system information: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const formatBytes = (bytes) => {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  const formatNumber = (num) => {
    return new Intl.NumberFormat().format(num);
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <InfoIcon sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Yadoma - System Information
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <IconButton color="inherit" onClick={loadData}>
              <RefreshIcon />
            </IconButton>
            {user?.role === 'ADMIN' && (
              <IconButton color="inherit" onClick={() => navigate('/admin/users')}>
                <AdminIcon />
              </IconButton>
            )}
            <IconButton color="inherit" onClick={() => navigate('/dashboard')}>
              <DashboardIcon />
            </IconButton>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <AccountCircleIcon />
              <Typography variant="body1">{user?.email}</Typography>
              {user?.role === 'ADMIN' && (
                <Chip label="ADMIN" color="error" size="small" />
              )}
            </Box>
            <IconButton color="inherit" onClick={handleLogout}>
              <LogoutIcon />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            {systemInfo && (
              <Paper sx={{ p: 3, mb: 3 }}>
                <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <ComputerIcon color="primary" />
                  System Information
                </Typography>

                <Grid container spacing={2} sx={{ mt: 1 }}>
                  <Grid item xs={12} md={6} lg={3}>
                    <Card>
                      <CardContent>
                        <Typography color="text.secondary" gutterBottom>
                          Operating System
                        </Typography>
                        <Typography variant="h6">
                          {systemInfo.operatingSystem}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {systemInfo.architecture}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card>
                      <CardContent>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <MemoryIcon color="primary" />
                          <Typography color="text.secondary" gutterBottom>
                            CPU
                          </Typography>
                        </Box>
                        <Typography variant="h6">
                          {systemInfo.nCpu} Cores
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {systemInfo.kernelVersion}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card>
                      <CardContent>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <MemoryIcon color="success" />
                          <Typography color="text.secondary" gutterBottom>
                            Memory
                          </Typography>
                        </Box>
                        <Typography variant="h6">
                          {formatBytes(systemInfo.memTotal)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Total RAM
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card>
                      <CardContent>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <StorageIcon color="warning" />
                          <Typography color="text.secondary" gutterBottom>
                            Docker Version
                          </Typography>
                        </Box>
                        <Typography variant="h6">
                          {systemInfo.serverVersion}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {systemInfo.driver} driver
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card>
                      <CardContent>
                        <Typography color="text.secondary" gutterBottom>
                          Total Containers
                        </Typography>
                        <Typography variant="h6">
                          {systemInfo.containers}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card sx={{ bgcolor: 'success.light' }}>
                      <CardContent>
                        <Typography color="text.secondary" gutterBottom>
                          Running
                        </Typography>
                        <Typography variant="h6">
                          {systemInfo.containersRunning}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card sx={{ bgcolor: 'warning.light' }}>
                      <CardContent>
                        <Typography color="text.secondary" gutterBottom>
                          Paused
                        </Typography>
                        <Typography variant="h6">
                          {systemInfo.containersPaused}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6} lg={3}>
                    <Card sx={{ bgcolor: 'error.light' }}>
                      <CardContent>
                        <Typography color="text.secondary" gutterBottom>
                          Stopped
                        </Typography>
                        <Typography variant="h6">
                          {systemInfo.containersStopped}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6}>
                    <Card>
                      <CardContent>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <ImageIcon color="primary" />
                          <Typography color="text.secondary" gutterBottom>
                            Images
                          </Typography>
                        </Box>
                        <Typography variant="h6">
                          {systemInfo.images}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Docker images stored
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>

                  <Grid item xs={12} md={6}>
                    <Card>
                      <CardContent>
                        <Typography color="text.secondary" gutterBottom>
                          System ID
                        </Typography>
                        <Typography variant="body2" fontFamily="monospace">
                          {systemInfo.id}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                </Grid>
              </Paper>
            )}

            {diskUsage && (
              <Paper sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <StorageIcon color="primary" />
                  Disk Usage
                </Typography>

                <Box sx={{ mb: 3 }}>
                  <Typography variant="h6" gutterBottom>
                    Total Layers Size: {formatBytes(diskUsage.layersSize)}
                  </Typography>
                </Box>

                {diskUsage.images && diskUsage.images.length > 0 && (
                  <Box sx={{ mb: 3 }}>
                    <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <ImageIcon />
                      Images ({diskUsage.images.length})
                    </Typography>
                    <TableContainer>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Repository Tags</TableCell>
                            <TableCell>ID</TableCell>
                            <TableCell align="right">Size</TableCell>
                            <TableCell align="right">Containers</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {diskUsage.images.map((image) => (
                            <TableRow key={image.id}>
                              <TableCell>
                                {image.repoTags && image.repoTags.length > 0 ? (
                                  image.repoTags.map((tag, idx) => (
                                    <Chip key={idx} label={tag} size="small" sx={{ mr: 0.5 }} />
                                  ))
                                ) : (
                                  <Typography variant="body2" color="text.secondary">none</Typography>
                                )}
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2" fontFamily="monospace">
                                  {image.id.substring(0, 12)}
                                </Typography>
                              </TableCell>
                              <TableCell align="right">{formatBytes(image.size)}</TableCell>
                              <TableCell align="right">{image.containers}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Box>
                )}

                {diskUsage.volumes && diskUsage.volumes.length > 0 && (
                  <Box sx={{ mb: 3 }}>
                    <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <FolderIcon />
                      Volumes ({diskUsage.volumes.length})
                    </Typography>
                    <TableContainer>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Name</TableCell>
                            <TableCell>Mount Point</TableCell>
                            <TableCell align="right">Size</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {diskUsage.volumes.map((volume, idx) => (
                            <TableRow key={idx}>
                              <TableCell>
                                <Typography variant="body2" fontFamily="monospace">
                                  {volume.name}
                                </Typography>
                              </TableCell>
                              <TableCell>{volume.mountpoint}</TableCell>
                              <TableCell align="right">{formatBytes(volume.size)}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Box>
                )}

                {diskUsage.containers && diskUsage.containers.length > 0 && (
                  <Box>
                    <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <ViewInArIcon />
                      Container Disk Usage ({diskUsage.containers.length})
                    </Typography>
                    <TableContainer>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Image</TableCell>
                            <TableCell>State</TableCell>
                            <TableCell>Status</TableCell>
                            <TableCell align="right">Size (RW)</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {diskUsage.containers.map((container) => (
                            <TableRow key={container.id}>
                              <TableCell>
                                <Typography variant="body2" fontFamily="monospace">
                                  {container.id.substring(0, 12)}
                                </Typography>
                              </TableCell>
                              <TableCell>{container.image}</TableCell>
                              <TableCell>
                                <Chip
                                  label={container.state}
                                  size="small"
                                  color={container.state === 'running' ? 'success' : 'default'}
                                />
                              </TableCell>
                              <TableCell>{container.status}</TableCell>
                              <TableCell align="right">{formatBytes(container.sizeRw)}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Box>
                )}
              </Paper>
            )}
          </>
        )}
      </Container>
    </Box>
  );
};

export default SystemInfoPage;
