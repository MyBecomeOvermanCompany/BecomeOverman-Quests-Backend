import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Card, 
  CardContent, 
  Typography, 
  Button, 
  Grid, 
  CircularProgress,
  Avatar,
  Stack,
  Chip,
  Paper,
} from '@mui/material';
import PersonIcon from '@mui/icons-material/Person';
import EmailIcon from '@mui/icons-material/Email';
import StarIcon from '@mui/icons-material/Star';
import BoltIcon from '@mui/icons-material/Bolt';
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';
import RefreshIcon from '@mui/icons-material/Refresh';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

function ProfileTab() {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadProfile();
  }, []);

  const loadProfile = async () => {
    setLoading(true);
    try {
      const result = await apiCall('/user/profile');
      
      if (result.success && result.data) {
        setProfile(result.data);
        showToast('Профиль загружен!', 'success');
      } else {
        showToast('Ошибка загрузки профиля: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки профиля', 'error');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress size={60} />
      </Box>
    );
  }

  if (!profile) {
    return (
      <Box sx={{ textAlign: 'center', py: 8 }}>
        <Typography variant="h6" color="text.secondary" gutterBottom>
          Не удалось загрузить профиль
        </Typography>
        <Button onClick={loadProfile} sx={{ mt: 2 }}>
          Повторить
        </Button>
      </Box>
    );
  }

  const createdAt = new Date(profile.created_at).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  return (
    <Box>
      <Typography variant="h3" gutterBottom sx={{ mb: 4, fontWeight: 700 }}>
        Профиль
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Stack spacing={3} alignItems="center">
                <Avatar 
                  sx={{ 
                    width: 120, 
                    height: 120, 
                    bgcolor: 'primary.main', 
                    fontSize: '3rem',
                    border: '4px solid',
                    borderColor: 'primary.main',
                  }}
                >
                  {profile.username?.charAt(0).toUpperCase() || 'U'}
                </Avatar>
                <Box textAlign="center">
                  <Typography variant="h5" gutterBottom fontWeight={600}>
                    {profile.username}
                  </Typography>
                  <Chip 
                    label="Активный игрок" 
                    color="success" 
                    size="small"
                    sx={{ mt: 1 }}
                  />
                </Box>
                <Button
                  variant="outlined"
                  startIcon={<RefreshIcon />}
                  onClick={loadProfile}
                  fullWidth
                >
                  Обновить
                </Button>
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
                Информация
              </Typography>
              <Grid container spacing={3}>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2, bgcolor: '#fafafa' }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <PersonIcon color="primary" />
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          ID пользователя
                        </Typography>
                        <Typography variant="h6" fontWeight={600}>
                          {profile.id}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2, bgcolor: '#fafafa' }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <EmailIcon color="primary" />
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Email
                        </Typography>
                        <Typography variant="h6" fontWeight={600}>
                          {profile.email || 'Не указан'}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2, bgcolor: '#fafafa' }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <StarIcon sx={{ color: '#fbbf24' }} />
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Уровень
                        </Typography>
                        <Typography variant="h6" fontWeight={600}>
                          {profile.level || 1}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2, bgcolor: '#fafafa' }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <BoltIcon sx={{ color: '#10b981' }} />
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Опыт
                        </Typography>
                        <Typography variant="h6" fontWeight={600}>
                          {profile.xp_points || 0} XP
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2, bgcolor: '#fafafa' }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <AccountBalanceWalletIcon sx={{ color: '#f59e0b' }} />
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Монеты
                        </Typography>
                        <Typography variant="h6" fontWeight={600}>
                          {profile.coin_balance || 0}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Paper sx={{ p: 2, bgcolor: '#fafafa' }}>
                    <Stack direction="row" spacing={2} alignItems="center">
                      <CalendarTodayIcon color="primary" />
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Регистрация
                        </Typography>
                        <Typography variant="h6" fontWeight={600}>
                          {createdAt}
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}

export default ProfileTab;
