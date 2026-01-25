import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardActions,
  Button,
  TextField,
  CircularProgress,
  Grid,
  Chip,
  Stack,
  Avatar,
  Paper,
} from '@mui/material';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import PersonIcon from '@mui/icons-material/Person';
import SearchIcon from '@mui/icons-material/Search';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import LocalFireDepartmentIcon from '@mui/icons-material/LocalFireDepartment';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import Divider from '@mui/material/Divider';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

function FriendsTab() {
  const [friends, setFriends] = useState([]);
  const [friendRequests, setFriendRequests] = useState([]);
  const [friendRecommendations, setFriendRecommendations] = useState([]);
  const [sharedQuests, setSharedQuests] = useState([]);
  const [availableQuests, setAvailableQuests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [friendName, setFriendName] = useState('');
  const [selectedFriendId, setSelectedFriendId] = useState('');
  const [selectedQuestId, setSelectedQuestId] = useState('');
  const [sharedQuestDialogOpen, setSharedQuestDialogOpen] = useState(false);
  const [friendStats, setFriendStats] = useState({});
  const [statsDialogOpen, setStatsDialogOpen] = useState(false);
  const [selectedFriendForStats, setSelectedFriendForStats] = useState(null);

  useEffect(() => {
    loadFriends();
    loadSharedQuests();
    loadFriendRequests();
    loadAvailableQuests();
  }, []);

  const loadFriends = async () => {
    setLoading(true);
    try {
      const result = await apiCall('/friends');

      if (result.success && result.data && Array.isArray(result.data)) {
        setFriends(result.data);
      }
    } catch (error) {
      showToast('Ошибка загрузки друзей', 'error');
    } finally {
      setLoading(false);
    }
  };

  const loadSharedQuests = async () => {
    try {
      const result = await apiCall('/quests/shared');
      if (result.success && result.data && Array.isArray(result.data)) {
        setSharedQuests(result.data);
      } else {
        setSharedQuests([]);
      }
    } catch (error) {
      console.error('Error loading shared quests:', error);
      setSharedQuests([]);
    }
  };

  const loadFriendRequests = async () => {
    try {
      const result = await apiCall('/friends/requests');
      if (result.success && result.data && Array.isArray(result.data)) {
        setFriendRequests(result.data);
      }
    } catch (error) {
      console.error('Error loading friend requests:', error);
    }
  };

  const loadAvailableQuests = async () => {
    try {
      const result = await apiCall('/quests/available');
      if (result.success && Array.isArray(result)) {
        setAvailableQuests(result);
      } else if (result.data && Array.isArray(result.data)) {
        setAvailableQuests(result.data);
      }
    } catch (error) {
      console.error('Error loading available quests:', error);
    }
  };

  const loadFriendStats = async (friendId) => {
    try {
      const result = await apiCall(`/friends/${friendId}/stats`);
      if (result.success && result.data) {
        setFriendStats(prev => ({ ...prev, [friendId]: result.data }));
        return result.data;
      }
    } catch (error) {
      console.error('Error loading friend stats:', error);
    }
    return null;
  };

  const addFriendByName = async () => {
    if (!friendName.trim()) {
      showToast('Введите имя друга', 'warning');
      return;
    }

    try {
      const result = await apiCall(`/friends/by-name/${friendName}`, {
        method: 'POST',
      });

      if (result.success || result.message) {
        showToast('Запрос в друзья отправлен!', 'success');
        setFriendName('');
        loadFriendRequests();
      } else {
        showToast('Ошибка: ' + (result.error || result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка отправки запроса', 'error');
    }
  };

  const addFriendById = async (friendId) => {
    try {
      const result = await apiCall(`/friends/${friendId}`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Друг добавлен!', 'success');
        loadFriends();
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка добавления друга', 'error');
    }
  };

  const loadFriendRecommendations = async () => {
    try {
      const result = await apiCall('/quests/recommend/friends', {
        method: 'POST',
      });

      if (result.success && result.data && Array.isArray(result.data)) {
        setFriendRecommendations(result.data);
        if (result.data.length === 0) {
          showToast('Рекомендации друзей пока недоступны', 'info');
        }
      } else {
        showToast('Ошибка загрузки рекомендаций', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки рекомендаций', 'error');
    }
  };

  const acceptFriendRequest = async (friendId) => {
    try {
      const result = await apiCall(`/friends/${friendId}/accept`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Запрос в друзья принят!', 'success');
        loadFriendRequests();
        loadFriends();
      } else {
        showToast('Ошибка: ' + (result.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка принятия запроса', 'error');
    }
  };

  const rejectFriendRequest = async (friendId) => {
    try {
      const result = await apiCall(`/friends/${friendId}/reject`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Запрос в друзья отклонен', 'info');
        loadFriendRequests();
      } else {
        showToast('Ошибка: ' + (result.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка отклонения запроса', 'error');
    }
  };

  const showFriendStats = async (friend) => {
    setSelectedFriendForStats(friend);
    setStatsDialogOpen(true);
    if (!friendStats[friend.friend_id]) {
      await loadFriendStats(friend.friend_id);
    }
  };

  const createSharedQuest = async () => {
    if (!selectedFriendId || !selectedQuestId) {
      showToast('Выберите друга и квест', 'warning');
      return;
    }

    try {
      const result = await apiCall('/quests/shared', {
        method: 'POST',
        body: JSON.stringify({ 
          friend_id: parseInt(selectedFriendId),
          quest_id: parseInt(selectedQuestId)
        }),
      });

      if (result.success || result.message) {
        showToast('Совместный квест создан!', 'success');
        setSelectedFriendId('');
        setSelectedQuestId('');
        setSharedQuestDialogOpen(false);
        loadSharedQuests();
        loadFriends();
      } else {
        showToast('Ошибка: ' + (result.error || result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка создания совместного квеста', 'error');
    }
  };

  return (
    <Box>
      <Typography variant="h3" gutterBottom sx={{ mb: 4, fontWeight: 700 }}>
        Друзья
      </Typography>

      <Paper 
        elevation={0}
        sx={{ 
          p: 3, 
          mb: 4,
          border: '1px solid',
          borderColor: 'divider',
        }}
      >
        <Stack spacing={2}>
          <Typography variant="h6">Добавить друга</Typography>
          <Box sx={{ display: 'flex', gap: 2 }}>
            <TextField
              fullWidth
              label="Имя пользователя"
              value={friendName}
              onChange={(e) => setFriendName(e.target.value)}
              InputProps={{
                startAdornment: (
                  <SearchIcon sx={{ mr: 1, color: 'text.secondary' }} />
                ),
              }}
            />
            <Button 
              variant="contained" 
              onClick={addFriendByName}
              startIcon={<PersonAddIcon />}
              sx={{ minWidth: 150 }}
            >
              Добавить
            </Button>
          </Box>
        </Stack>
      </Paper>

      {/* Запросы в друзья */}
      {friendRequests.length > 0 && (
        <Box sx={{ mb: 4 }}>
          <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>
            Запросы в друзья ({friendRequests.length})
          </Typography>
          <Grid container spacing={2}>
            {friendRequests.map(request => (
              <Grid item xs={12} sm={6} md={4} key={request.id}>
                <Card>
                  <CardContent>
                    <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 2 }}>
                      <Avatar sx={{ bgcolor: 'warning.main' }}>
                        {request.username?.charAt(0).toUpperCase() || 'U'}
                      </Avatar>
                      <Box sx={{ flexGrow: 1 }}>
                        <Typography variant="subtitle1" fontWeight={600}>
                          {request.username || 'Без имени'}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Хочет добавить вас в друзья
                        </Typography>
                      </Box>
                    </Stack>
                  </CardContent>
                  <CardActions>
                    <Button
                      size="small"
                      variant="contained"
                      color="success"
                      startIcon={<CheckCircleIcon />}
                      onClick={() => acceptFriendRequest(request.friend_id)}
                      sx={{ flex: 1 }}
                    >
                      Принять
                    </Button>
                    <Button
                      size="small"
                      variant="outlined"
                      color="error"
                      startIcon={<CancelIcon />}
                      onClick={() => rejectFriendRequest(request.friend_id)}
                      sx={{ flex: 1 }}
                    >
                      Отклонить
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Box>
      )}

      <Typography variant="h5" sx={{ mt: 4, mb: 2, fontWeight: 600 }}>
        Мои друзья
      </Typography>

      {loading ? (
        <Box display="flex" justifyContent="center" py={4}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2}>
          {friends.map(friend => (
            <Grid item xs={12} sm={6} md={4} key={friend.id}>
              <Card>
                <CardContent>
                  <Stack direction="row" spacing={2} alignItems="center">
                    <Avatar sx={{ bgcolor: 'primary.main' }}>
                      {friend.username?.charAt(0).toUpperCase() || 'F'}
                    </Avatar>
                    <Box sx={{ flexGrow: 1 }}>
                      <Typography variant="subtitle1" fontWeight={600}>
                        {friend.username || 'Без имени'}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        ID: {friend.friend_id}
                      </Typography>
                    </Box>
                  </Stack>
                </CardContent>
                <CardActions>
                  <Button
                    size="small"
                    variant="outlined"
                    onClick={() => showFriendStats(friend)}
                    fullWidth
                  >
                    Статистика
                  </Button>
                  <Button
                    size="small"
                    variant="contained"
                    onClick={() => {
                      setSelectedFriendId(friend.friend_id.toString());
                      setSharedQuestDialogOpen(true);
                    }}
                    fullWidth
                  >
                    Совместный квест
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {friends.length === 0 && !loading && (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            У вас пока нет друзей
          </Typography>
        </Box>
      )}

      <Typography variant="h5" sx={{ mt: 4, mb: 2, fontWeight: 600 }}>
        Рекомендации друзей
      </Typography>
      <Button 
        variant="contained" 
        onClick={loadFriendRecommendations} 
        fullWidth 
        sx={{ mb: 3 }}
        startIcon={<PersonAddIcon />}
      >
        Найти друзей по интересам
      </Button>

      {friendRecommendations.length > 0 && (
        <Grid container spacing={2}>
          {friendRecommendations.map((friend, idx) => (
            <Grid item xs={12} sm={6} md={4} key={friend.id || idx}>
              <Card>
                <CardContent>
                  <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 2 }}>
                    <Avatar sx={{ bgcolor: 'secondary.main' }}>
                      {friend.username?.charAt(0).toUpperCase() || 'U'}
                    </Avatar>
                    <Box sx={{ flexGrow: 1 }}>
                      <Typography variant="subtitle1" fontWeight={600}>
                        {friend.username || 'Без имени'}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        ID: {friend.id}
                      </Typography>
                    </Box>
                  </Stack>
                  {friend.similarity_score !== undefined && (
                    <Chip 
                      label={`${(friend.similarity_score * 100).toFixed(1)}% совпадение`} 
                      color="warning" 
                      size="small"
                      variant="outlined"
                    />
                  )}
                </CardContent>
                <CardActions>
                  <Button 
                    size="small" 
                    variant="contained" 
                    onClick={() => addFriendById(friend.id)}
                    fullWidth
                  >
                    Добавить в друзья
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {/* Диалог создания совместного квеста */}
      <Dialog 
        open={sharedQuestDialogOpen} 
        onClose={() => setSharedQuestDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Создать совместный квест</DialogTitle>
        <DialogContent>
          <Stack spacing={3} sx={{ mt: 1 }}>
            <FormControl fullWidth>
              <InputLabel>Выберите друга</InputLabel>
              <Select
                value={selectedFriendId}
                label="Выберите друга"
                onChange={(e) => setSelectedFriendId(e.target.value)}
              >
                {friends.map(friend => (
                  <MenuItem key={friend.friend_id} value={friend.friend_id.toString()}>
                    {friend.username} (ID: {friend.friend_id})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <FormControl fullWidth>
              <InputLabel>Выберите квест</InputLabel>
              <Select
                value={selectedQuestId}
                label="Выберите квест"
                onChange={(e) => setSelectedQuestId(e.target.value)}
              >
                {availableQuests.map(quest => (
                  <MenuItem key={quest.id} value={quest.id.toString()}>
                    {quest.title} ({quest.price} монет)
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSharedQuestDialogOpen(false)}>Отмена</Button>
          <Button 
            variant="contained" 
            color="success" 
            onClick={createSharedQuest}
            disabled={!selectedFriendId || !selectedQuestId}
          >
            Создать
          </Button>
        </DialogActions>
      </Dialog>

      {/* Диалог статистики друга */}
      <Dialog 
        open={statsDialogOpen} 
        onClose={() => setStatsDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Статистика: {selectedFriendForStats?.username || 'Друг'}
        </DialogTitle>
        <DialogContent>
          {selectedFriendForStats && friendStats[selectedFriendForStats.friend_id] ? (
            <Stack spacing={2} sx={{ mt: 1 }}>
              <Box>
                <Typography variant="body2" color="text.secondary">Уровень</Typography>
                <Typography variant="h6">{friendStats[selectedFriendForStats.friend_id].level || 1}</Typography>
              </Box>
              <Divider />
              <Box>
                <Typography variant="body2" color="text.secondary">Завершено квестов</Typography>
                <Typography variant="h6">
                  {friendStats[selectedFriendForStats.friend_id].total_quests_completed || 0}
                </Typography>
              </Box>
              <Divider />
              <Box>
                <Typography variant="body2" color="text.secondary">Выполнено задач</Typography>
                <Typography variant="h6">
                  {friendStats[selectedFriendForStats.friend_id].total_tasks_completed || 0}
                </Typography>
              </Box>
              <Divider />
              <Box>
                <Typography variant="body2" color="text.secondary">Текущий streak</Typography>
                <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <LocalFireDepartmentIcon color="error" />
                  {friendStats[selectedFriendForStats.friend_id].current_streak || 0} дней
                </Typography>
              </Box>
              <Divider />
              <Box>
                <Typography variant="body2" color="text.secondary">Лучший streak</Typography>
                <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <LocalFireDepartmentIcon color="warning" />
                  {friendStats[selectedFriendForStats.friend_id].longest_streak || 0} дней
                </Typography>
              </Box>
            </Stack>
          ) : (
            <Box sx={{ textAlign: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setStatsDialogOpen(false)}>Закрыть</Button>
        </DialogActions>
      </Dialog>

      {sharedQuests.length > 0 && (
        <Box sx={{ mt: 4 }}>
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
            Совместные квесты
          </Typography>
          <Grid container spacing={2}>
            {sharedQuests.map(quest => (
              <Grid item xs={12} sm={6} md={4} key={quest.id}>
                <Card>
                  <CardContent>
                    <Typography variant="subtitle1" fontWeight={600} gutterBottom>
                      {quest.title || 'Без названия'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                      {quest.description || 'Описание отсутствует'}
                    </Typography>
                    <Chip 
                      label={`С ${quest.friend_username || 'другом'}`} 
                      size="small" 
                      color="primary"
                      variant="outlined"
                    />
                    <Typography variant="caption" display="block" sx={{ mt: 1, color: 'text.secondary' }}>
                      Статус: {quest.status === 'active' ? 'Активен' : quest.status}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Box>
      )}
    </Box>
  );
}

export default FriendsTab;
