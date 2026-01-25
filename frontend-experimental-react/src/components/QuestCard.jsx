import React from 'react';
import { Card, CardContent, Box, Typography, Chip, Stack, LinearProgress, Button } from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import InfoIcon from '@mui/icons-material/Info';

function QuestCard({ quest, onAction, actionLabel = 'Начать', onDetails, showProgress = false, completed = false }) {
  if (!quest) return null;

  const progress = quest.tasks_count 
    ? ((quest.completed_tasks_count || 0) / quest.tasks_count) * 100 
    : 0;

  const getInitials = (title) => {
    if (!title) return 'Q';
    const words = title.split(' ').filter(w => w.length > 0);
    if (words.length >= 2) {
      return (words[0][0] + words[1][0]).toUpperCase();
    }
    return title.substring(0, 2).toUpperCase();
  };

  const handleClick = () => {
    if (onAction) {
      onAction(quest);
    }
  };

  const handleDetailsClick = () => {
    if (onDetails) {
      onDetails(quest);
    }
  };

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box
        sx={{
          height: 180,
          background: `linear-gradient(135deg, #000000 0%, #1a1a1a 100%)`,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          position: 'relative',
        }}
      >
        <Typography
          variant="h2"
          sx={{
            color: '#ffffff',
            fontWeight: 700,
            fontSize: '3rem',
            letterSpacing: '-0.02em',
          }}
        >
          {getInitials(quest.title || 'Quest')}
        </Typography>
      </Box>
      
      <CardContent sx={{ p: 3, flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
        <Typography 
          variant="h6" 
          sx={{ 
            mb: 1.5, 
            fontWeight: 600,
            lineHeight: 1.3,
            minHeight: '3em',
          }}
        >
          {quest.title || 'Без названия'}
        </Typography>
        
        <Typography 
          variant="body2" 
          color="text.secondary"
          sx={{ 
            mb: 2, 
            flexGrow: 1,
            display: '-webkit-box',
            WebkitLineClamp: 4,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
            minHeight: '5em',
            lineHeight: 1.5,
          }}
        >
          {quest.description || 'Описание отсутствует'}
        </Typography>
        
        <Stack direction="row" spacing={1} sx={{ mb: 2 }} flexWrap="wrap">
          {quest.price !== undefined && quest.price !== null && (
            <Chip 
              label={`${quest.price} монет`} 
              size="small"
              sx={{ 
                bgcolor: '#fee2e2', 
                color: '#991b1b',
                fontWeight: 600,
                border: '1px solid #fca5a5'
              }}
            />
          )}
          {quest.category && (
            <Chip 
              label={quest.category} 
              size="small" 
              sx={{ bgcolor: '#f3f4f6' }}
            />
          )}
          {quest.reward_xp && (
            <Chip 
              label={`${quest.reward_xp} XP`} 
              size="small"
              sx={{ bgcolor: '#dbeafe', color: '#1e40af' }}
            />
          )}
          {quest.reward_coin && (
            <Chip 
              label={`+${quest.reward_coin} монет`} 
              size="small"
              sx={{ bgcolor: '#fef3c7', color: '#92400e' }}
            />
          )}
          {completed && (
            <Chip 
              icon={<CheckCircleIcon />}
              label="Завершен" 
              size="small"
              color="success"
            />
          )}
        </Stack>

        {showProgress && quest.tasks_count && (
          <Box sx={{ mb: 2 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
              <Typography variant="caption" color="text.secondary">
                Прогресс
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {quest.completed_tasks_count || 0} / {quest.tasks_count}
              </Typography>
            </Box>
            <LinearProgress 
              variant="determinate" 
              value={progress}
              sx={{ height: 6, borderRadius: 3, bgcolor: '#f3f4f6' }}
            />
          </Box>
        )}

        <Stack direction="row" spacing={1}>
          {onDetails && (
            <Button
              variant="outlined"
              onClick={handleDetailsClick}
              startIcon={<InfoIcon />}
              sx={{ flex: 1 }}
            >
              Детали
            </Button>
          )}
          {onAction && (
            <Button
              variant="contained"
              onClick={handleClick}
              startIcon={<PlayArrowIcon />}
              sx={{ flex: 1 }}
            >
              {actionLabel}
            </Button>
          )}
        </Stack>
      </CardContent>
    </Card>
  );
}

export default QuestCard;
