import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Box,
  Button,
  AppBar,
  Toolbar,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Alert,
  CircularProgress,
  Tooltip,
  Collapse,
} from '@mui/material';
import {
  Logout as LogoutIcon,
  Delete as DeleteIcon,
  AccountCircle as AccountCircleIcon,
  Dashboard as DashboardIcon,
  AdminPanelSettings as AdminIcon,
  KeyboardArrowDown as KeyboardArrowDownIcon,
  KeyboardArrowUp as KeyboardArrowUpIcon,
  PlayArrow as PlayArrowIcon,
  Stop as StopIcon,
  Refresh as RefreshIcon,
  Storage as StorageIcon,
} from '@mui/icons-material';
import { useAuth } from '../context/AuthContext';
import { userApi, containerApi } from '../services/api';

const AdminUsersPage = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState({});
  const [expandedUser, setExpandedUser] = useState(null);
  const [userContainers, setUserContainers] = useState({});

  useEffect(() => {
    if (user) {
      loadUsers();
    }
  }, [user]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await userApi.getAllUsers();
      setUsers(data || []);
    } catch (err) {
      setError('Failed to load users: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleDeleteUser = async (userId, userEmail) => {
    if (!window.confirm(`Are you sure you want to delete user "${userEmail}"?`)) {
      return;
    }

    try {
      setActionLoading({ [userId]: 'delete' });
      await userApi.deleteUser(userId);
      await loadUsers();
    } catch (err) {
      setError('Failed to delete user: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleExpandUser = async (userId) => {
    if (expandedUser === userId) {
      setExpandedUser(null);
      return;
    }

    setExpandedUser(userId);

    if (!userContainers[userId]) {
      try {
        setActionLoading({ [userId]: 'loading-containers' });
        const containers = await containerApi.getAllContainers(userId);
        setUserContainers({ ...userContainers, [userId]: containers || [] });
      } catch (err) {
        setError('Failed to load containers: ' + (err.response?.data?.message || err.message));
      } finally {
        setActionLoading({});
      }
    }
  };

  const reloadUserContainers = async (userId) => {
    try {
      setActionLoading({ [userId]: 'loading-containers' });
      const containers = await containerApi.getAllContainers(userId);
      setUserContainers({ ...userContainers, [userId]: containers || [] });
    } catch (err) {
      setError('Failed to reload containers: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleStartContainer = async (userId, containerId) => {
    try {
      setActionLoading({ [containerId]: 'start' });
      await containerApi.startContainer(containerId);
      await reloadUserContainers(userId);
    } catch (err) {
      setError('Failed to start container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleStopContainer = async (userId, containerId) => {
    try {
      setActionLoading({ [containerId]: 'stop' });
      await containerApi.stopContainer(containerId);
      await reloadUserContainers(userId);
    } catch (err) {
      setError('Failed to stop container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleRestartContainer = async (userId, containerId) => {
    try {
      setActionLoading({ [containerId]: 'restart' });
      await containerApi.restartContainer(containerId);
      await reloadUserContainers(userId);
    } catch (err) {
      setError('Failed to restart container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleDeleteContainer = async (userId, containerId, containerName) => {
    if (!window.confirm(`Are you sure you want to delete container "${containerName}"?`)) {
      return;
    }

    try {
      setActionLoading({ [containerId]: 'delete' });
      await containerApi.deleteContainer(containerId);
      await reloadUserContainers(userId);
    } catch (err) {
      setError('Failed to delete container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const getRoleColor = (role) => {
    return role === 'ADMIN' ? 'error' : 'primary';
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'running':
        return 'success';
      case 'stopped':
        return 'error';
      case 'paused':
        return 'warning';
      default:
        return 'default';
    }
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <AdminIcon sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Yadoma - Admin Panel
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Button
              color="inherit"
              startIcon={<DashboardIcon />}
              onClick={() => navigate('/dashboard')}
            >
              Dashboard
            </Button>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <AccountCircleIcon />
              <Typography variant="body1">{user?.email}</Typography>
              <Chip label={user?.role} color="error" size="small" />
            </Box>
            <IconButton color="inherit" onClick={handleLogout}>
              <LogoutIcon />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Paper sx={{ p: 3 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h5" component="h2">
              User Management
            </Typography>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
              {error}
            </Alert>
          )}

          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : users.length === 0 ? (
            <Box sx={{ textAlign: 'center', py: 4 }}>
              <Typography variant="body1" color="text.secondary">
                No users found.
              </Typography>
            </Box>
          ) : (
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell width="50px" />
                    <TableCell>Email</TableCell>
                    <TableCell>Role</TableCell>
                    <TableCell>User ID</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {users.map((userItem) => (
                    <React.Fragment key={userItem.id}>
                      <TableRow>
                        <TableCell>
                          <IconButton
                            size="small"
                            onClick={() => handleExpandUser(userItem.id)}
                          >
                            {expandedUser === userItem.id ? (
                              <KeyboardArrowUpIcon />
                            ) : (
                              <KeyboardArrowDownIcon />
                            )}
                          </IconButton>
                        </TableCell>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <AccountCircleIcon color="action" />
                            {userItem.email}
                            {userItem.id === user?.id && (
                              <Chip label="You" size="small" color="primary" variant="outlined" />
                            )}
                          </Box>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={userItem.role}
                            color={getRoleColor(userItem.role)}
                            size="small"
                          />
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" fontFamily="monospace" color="text.secondary">
                            {userItem.id}
                          </Typography>
                        </TableCell>
                        <TableCell align="right">
                          <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
                            <Tooltip title={userItem.id === user?.id ? "You cannot delete yourself" : "Delete user"}>
                              <span>
                                <IconButton
                                  color="error"
                                  size="small"
                                  onClick={() => handleDeleteUser(userItem.id, userItem.email)}
                                  disabled={userItem.id === user?.id || !!actionLoading[userItem.id]}
                                >
                                  {actionLoading[userItem.id] === 'delete' ? (
                                    <CircularProgress size={20} />
                                  ) : (
                                    <DeleteIcon />
                                  )}
                                </IconButton>
                              </span>
                            </Tooltip>
                          </Box>
                        </TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={5}>
                          <Collapse in={expandedUser === userItem.id} timeout="auto" unmountOnExit>
                            <Box sx={{ margin: 2 }}>
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                                <StorageIcon color="primary" />
                                <Typography variant="h6" gutterBottom component="div">
                                  Containers
                                </Typography>
                              </Box>
                              {actionLoading[userItem.id] === 'loading-containers' ? (
                                <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
                                  <CircularProgress size={30} />
                                </Box>
                              ) : !userContainers[userItem.id] || userContainers[userItem.id].length === 0 ? (
                                <Typography variant="body2" color="text.secondary">
                                  No containers found for this user.
                                </Typography>
                              ) : (
                                <Table size="small">
                                  <TableHead>
                                    <TableRow>
                                      <TableCell>Name</TableCell>
                                      <TableCell>Status</TableCell>
                                      <TableCell>State</TableCell>
                                      <TableCell>Created At</TableCell>
                                      <TableCell align="right">Actions</TableCell>
                                    </TableRow>
                                  </TableHead>
                                  <TableBody>
                                    {userContainers[userItem.id].map((container) => (
                                      <TableRow key={container.id}>
                                        <TableCell>{container.name}</TableCell>
                                        <TableCell>
                                          <Chip
                                            label={container.status}
                                            color={getStatusColor(container.status)}
                                            size="small"
                                          />
                                        </TableCell>
                                        <TableCell>{container.state}</TableCell>
                                        <TableCell>
                                          {container.createdAt
                                            ? new Date(container.createdAt).toLocaleString()
                                            : 'N/A'}
                                        </TableCell>
                                        <TableCell align="right">
                                          <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
                                            <Tooltip title="Start">
                                              <IconButton
                                                color="success"
                                                size="small"
                                                onClick={() => handleStartContainer(userItem.id, container.id)}
                                                disabled={!!actionLoading[container.id]}
                                              >
                                                {actionLoading[container.id] === 'start' ? (
                                                  <CircularProgress size={16} />
                                                ) : (
                                                  <PlayArrowIcon fontSize="small" />
                                                )}
                                              </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Stop">
                                              <IconButton
                                                color="error"
                                                size="small"
                                                onClick={() => handleStopContainer(userItem.id, container.id)}
                                                disabled={!!actionLoading[container.id]}
                                              >
                                                {actionLoading[container.id] === 'stop' ? (
                                                  <CircularProgress size={16} />
                                                ) : (
                                                  <StopIcon fontSize="small" />
                                                )}
                                              </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Restart">
                                              <IconButton
                                                color="primary"
                                                size="small"
                                                onClick={() => handleRestartContainer(userItem.id, container.id)}
                                                disabled={!!actionLoading[container.id]}
                                              >
                                                {actionLoading[container.id] === 'restart' ? (
                                                  <CircularProgress size={16} />
                                                ) : (
                                                  <RefreshIcon fontSize="small" />
                                                )}
                                              </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Delete">
                                              <IconButton
                                                color="error"
                                                size="small"
                                                onClick={() => handleDeleteContainer(userItem.id, container.id, container.name)}
                                                disabled={!!actionLoading[container.id]}
                                              >
                                                {actionLoading[container.id] === 'delete' ? (
                                                  <CircularProgress size={16} />
                                                ) : (
                                                  <DeleteIcon fontSize="small" />
                                                )}
                                              </IconButton>
                                            </Tooltip>
                                          </Box>
                                        </TableCell>
                                      </TableRow>
                                    ))}
                                  </TableBody>
                                </Table>
                              )}
                            </Box>
                          </Collapse>
                        </TableCell>
                      </TableRow>
                    </React.Fragment>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Paper>
      </Container>
    </Box>
  );
};

export default AdminUsersPage;
