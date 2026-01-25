import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  TextField,
  InputAdornment,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Chip,
  Stack,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import StarIcon from '@mui/icons-material/Star';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import Paper from '@mui/material/Paper';
import Divider from '@mui/material/Divider';
import { apiCall } from '../api';
import QuestCard from './QuestCard';
import MasonryGrid from './MasonryGrid';
import { showToast } from '../utils/toast';

function QuestsTab() {
  const [availableQuests, setAvailableQuests] = useState([]);
  const [shopQuests, setShopQuests] = useState([]);
  const [recommendations, setRecommendations] = useState([]);
  const [questDetails, setQuestDetails] = useState(null);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  const [searchAvailable, setSearchAvailable] = useState('');
  const [searchShop, setSearchShop] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadQuests();
  }, []);

  const loadQuests = async () => {
    setLoading(true);
    try {
      const available = await apiCall('/quests/available');
      const shop = await apiCall('/quests/shop');

      if (available.success && available.data && Array.isArray(available.data)) {
        setAvailableQuests(available.data);
      }

      if (shop.success && shop.data && Array.isArray(shop.data)) {
        setShopQuests(shop.data);
      }
    } catch (error) {
      showToast('Ошибка загрузки квестов', 'error');
    } finally {
      setLoading(false);
    }
  };

  const loadRecommendations = async () => {
    try {
      const result = await apiCall('/quests/recommend', {
        method: 'POST',
      });

      if (result.success && result.data && result.data.recommendations && Array.isArray(result.data.recommendations)) {
        setRecommendations(result.data.recommendations);
        if (result.data.recommendations.length === 0) {
          showToast('Рекомендации пока недоступны. Попробуйте выполнить несколько квестов!', 'info');
        }
      } else {
        showToast('Ошибка загрузки рекомендаций', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки рекомендаций', 'error');
    }
  };

  const purchaseQuest = async (quest) => {
    if (!quest || !quest.id) return;
    
    try {
      const result = await apiCall(`/quests/${quest.id}/purchase`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Квест успешно куплен!', 'success');
        loadQuests();
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка покупки квеста', 'error');
    }
  };

  const showQuestDetails = async (quest) => {
    if (!quest || !quest.id) return;
    
    try {
      const result = await apiCall(`/quests/${quest.id}/details`);
      if (result.success) {
        setQuestDetails(result.data);
        setDetailsDialogOpen(true);
      } else {
        showToast('Ошибка загрузки деталей', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки деталей', 'error');
    }
  };

  const filterQuests = (quests, searchTerm) => {
    if (!searchTerm) return quests;
    if (!Array.isArray(quests)) return [];
    const term = searchTerm.toLowerCase();
    return quests.filter(quest => 
      (quest.title && quest.title.toLowerCase().includes(term)) ||
      (quest.description && quest.description.toLowerCase().includes(term))
    );
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress size={60} />
      </Box>
    );
  }

  const filteredAvailable = filterQuests(availableQuests, searchAvailable);
  const filteredShop = filterQuests(shopQuests, searchShop);

  return (
    <Box sx={{ pb: 4 }}>
      <Typography variant="h3" gutterBottom sx={{ mb: 4, fontWeight: 700 }}>
        Квесты
      </Typography>

      {filteredAvailable.length > 0 && (
        <Box sx={{ mb: 6 }}>
          <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
            Доступные квесты
          </Typography>
          <TextField
            fullWidth
            placeholder="Поиск..."
            value={searchAvailable}
            onChange={(e) => setSearchAvailable(e.target.value)}
            sx={{ mb: 3, maxWidth: 400 }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
          />
          <MasonryGrid>
            {filteredAvailable.map(quest => (
              <QuestCard
                key={quest.id}
                quest={quest}
                onAction={purchaseQuest}
                onDetails={showQuestDetails}
                actionLabel="Купить"
              />
            ))}
          </MasonryGrid>
        </Box>
      )}

      {filteredShop.length > 0 && (
        <Box sx={{ mb: 6 }}>
          <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
            Магазин квестов
          </Typography>
          <TextField
            fullWidth
            placeholder="Поиск..."
            value={searchShop}
            onChange={(e) => setSearchShop(e.target.value)}
            sx={{ mb: 3, maxWidth: 400 }}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
          />
          <MasonryGrid>
            {filteredShop.map(quest => (
              <QuestCard
                key={quest.id}
                quest={quest}
                onAction={purchaseQuest}
                onDetails={showQuestDetails}
                actionLabel="Купить"
              />
            ))}
          </MasonryGrid>
        </Box>
      )}

      <Box sx={{ mb: 6 }}>
        <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
          Рекомендации для вас
        </Typography>
        <Button 
          variant="contained" 
          onClick={loadRecommendations} 
          sx={{ mb: 3 }}
        >
          Получить персональные рекомендации
        </Button>
        {recommendations.length > 0 && (
          <MasonryGrid>
            {recommendations.map((quest, idx) => (
              <QuestCard
                key={quest.id || idx}
                quest={quest}
                onDetails={showQuestDetails}
                onAction={purchaseQuest}
                actionLabel="Купить"
              />
            ))}
          </MasonryGrid>
        )}
      </Box>

      <Dialog
        open={detailsDialogOpen}
        onClose={() => setDetailsDialogOpen(false)}
        maxWidth="md"
        fullWidth
        PaperProps={{
          sx: {
            borderRadius: 3,
          }
        }}
      >
        <DialogTitle 
          sx={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            pb: 2,
            background: 'linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%)',
            color: 'white',
            borderRadius: '12px 12px 0 0',
          }}
        >
          <Typography variant="h4" fontWeight={700} sx={{ flex: 1 }}>
            {questDetails?.title || 'Детали квеста'}
          </Typography>
          <IconButton 
            onClick={() => setDetailsDialogOpen(false)}
            sx={{ color: 'white' }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent sx={{ p: 0 }}>
          {questDetails && (
            <Box sx={{ p: 3 }}>
              <Stack spacing={3}>
                {/* Описание */}
                <Paper 
                  elevation={0}
                  sx={{ 
                    p: 3, 
                    bgcolor: '#f8f9fa',
                    borderRadius: 2,
                    border: '1px solid #e9ecef'
                  }}
                >
                  <Typography 
                    variant="body1" 
                    sx={{ 
                      lineHeight: 1.8,
                      fontSize: '1.05rem',
                      color: '#495057'
                    }}
                  >
                    {questDetails.description}
                  </Typography>
                </Paper>

                {/* Награды */}
                <Box>
                  <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                    Награды
                  </Typography>
                  <Stack direction="row" spacing={2} flexWrap="wrap">
                    <Chip 
                      icon={<AttachMoneyIcon />}
                      label={`Цена: ${questDetails.price || 0} монет`}
                      sx={{ 
                        bgcolor: '#fee2e2',
                        color: '#991b1b',
                        fontWeight: 600,
                        fontSize: '0.95rem',
                        height: 36,
                        '& .MuiChip-icon': {
                          color: '#991b1b'
                        }
                      }}
                    />
                    <Chip 
                      icon={<StarIcon />}
                      label={`Награда: ${questDetails.reward_coin || 0} монет`}
                      sx={{ 
                        bgcolor: '#fef3c7',
                        color: '#92400e',
                        fontWeight: 600,
                        fontSize: '0.95rem',
                        height: 36,
                        '& .MuiChip-icon': {
                          color: '#92400e'
                        }
                      }}
                    />
                    <Chip 
                      icon={<TrendingUpIcon />}
                      label={`${questDetails.reward_xp || 0} XP`}
                      sx={{ 
                        bgcolor: '#dbeafe',
                        color: '#1e40af',
                        fontWeight: 600,
                        fontSize: '0.95rem',
                        height: 36,
                        '& .MuiChip-icon': {
                          color: '#1e40af'
                        }
                      }}
                    />
                  </Stack>
                </Box>

                <Divider />

                {/* Задачи */}
                {questDetails.tasks && questDetails.tasks.length > 0 && (
                  <Box>
                    <Typography variant="h6" gutterBottom sx={{ mb: 2, fontWeight: 600 }}>
                      Задачи ({questDetails.tasks.length})
                    </Typography>
                    <Stack spacing={2}>
                      {questDetails.tasks.map((task, idx) => (
                        <Paper
                          key={idx}
                          elevation={0}
                          sx={{
                            p: 2.5,
                            border: '1px solid #e9ecef',
                            borderRadius: 2,
                            bgcolor: '#ffffff',
                            transition: 'all 0.2s',
                            '&:hover': {
                              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                              borderColor: '#4a5568',
                            }
                          }}
                        >
                          <Stack direction="row" spacing={2} alignItems="flex-start">
                            <Box
                              sx={{
                                minWidth: 32,
                                height: 32,
                                borderRadius: '50%',
                                bgcolor: '#2d3748',
                                color: 'white',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                fontWeight: 700,
                                fontSize: '0.9rem',
                                flexShrink: 0,
                              }}
                            >
                              {idx + 1}
                            </Box>
                            <Box sx={{ flex: 1 }}>
                              <Typography 
                                variant="subtitle1" 
                                fontWeight={600}
                                gutterBottom
                                sx={{ mb: 1 }}
                              >
                                {task.title}
                              </Typography>
                              <Typography 
                                variant="body2" 
                                color="text.secondary"
                                sx={{ lineHeight: 1.6 }}
                              >
                                {task.description}
                              </Typography>
                            </Box>
                          </Stack>
                        </Paper>
                      ))}
                    </Stack>
                  </Box>
                )}
              </Stack>
            </Box>
          )}
        </DialogContent>
        <DialogActions sx={{ p: 2.5, bgcolor: '#f8f9fa', borderRadius: '0 0 12px 12px' }}>
          <Button 
            onClick={() => setDetailsDialogOpen(false)}
            sx={{ minWidth: 100 }}
          >
            Закрыть
          </Button>
          {questDetails && (
            <Button 
              variant="contained" 
              onClick={() => {
                purchaseQuest(questDetails);
                setDetailsDialogOpen(false);
              }}
              startIcon={<AttachMoneyIcon />}
              sx={{ 
                minWidth: 150,
                bgcolor: '#2d3748',
                '&:hover': {
                  bgcolor: '#1a202c',
                }
              }}
            >
              Купить квест
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default QuestsTab;
