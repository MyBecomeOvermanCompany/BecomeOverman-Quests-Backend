import React, { useState, useEffect } from 'react';
import { ThemeProvider, CssBaseline } from '@mui/material';
import { Box, AppBar, Toolbar, Typography, Button, Container } from '@mui/material';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { isAuthenticated, logout } from './api';
import AuthTab from './components/AuthTab';
import ProfileTab from './components/ProfileTab';
import QuestsTab from './components/QuestsTab';
import MyQuestsTab from './components/MyQuestsTab';
import FriendsTab from './components/FriendsTab';
import AIQuestsTab from './components/AIQuestsTab';
import SearchQuestsTab from './components/SearchQuestsTab';
import CalendarTab from './components/CalendarTab';
import QuestMapTab from './components/QuestMapTab';
import RoadmapTab from './components/RoadmapTab';
import theme from './theme';
import './index.css';

function App() {
  const [activeTab, setActiveTab] = useState('quests');
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    setAuthenticated(isAuthenticated());
    if (isAuthenticated()) {
      setActiveTab('quests');
    } else {
      setActiveTab('auth');
    }
  }, []);

  const handleLogin = () => {
    setAuthenticated(true);
    setActiveTab('quests');
  };

  const handleLogout = () => {
    logout();
    setAuthenticated(false);
    setActiveTab('auth');
  };

  const tabs = [
    { id: 'auth', label: 'Авторизация', show: !authenticated },
    { id: 'profile', label: 'Профиль', show: authenticated },
    { id: 'roadmap', label: 'Roadmap', show: authenticated },
    { id: 'quests', label: 'Квесты', show: authenticated },
    { id: 'my-quests', label: 'Мои квесты', show: authenticated },
    { id: 'ai-quests', label: 'AI Квесты', show: authenticated },
    { id: 'search-quests', label: 'Поиск', show: authenticated },
    { id: 'friends', label: 'Друзья', show: authenticated },
    { id: 'calendar', label: 'Календарь', show: authenticated },
    { id: 'quest-map', label: 'Карта', show: authenticated },
  ].filter(tab => tab.show);

  const getTabContent = () => {
    switch (activeTab) {
      case 'auth':
        return <AuthTab onLogin={handleLogin} />;
      case 'profile':
        return <ProfileTab />;
      case 'roadmap':
        return <RoadmapTab />;
      case 'quests':
        return <QuestsTab />;
      case 'ai-quests':
        return <AIQuestsTab />;
      case 'search-quests':
        return <SearchQuestsTab />;
      case 'friends':
        return <FriendsTab />;
      case 'my-quests':
        return <MyQuestsTab />;
      case 'calendar':
        return <CalendarTab />;
      case 'quest-map':
        return <QuestMapTab />;
      default:
        return null;
    }
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ minHeight: '100vh', bgcolor: 'background.default' }}>
        <AppBar position="sticky" color="default">
          <Toolbar sx={{ justifyContent: 'space-between' }}>
            <Typography 
              variant="h5" 
              component="div" 
              sx={{ 
                fontWeight: 700,
                color: 'primary.main',
              }}
            >
              Become OverMan
            </Typography>
            {authenticated && (
              <Box sx={{ display: 'flex', gap: 1, alignItems: 'center', flexWrap: 'wrap' }}>
                {tabs.filter(t => t.id !== 'auth').map(tab => (
                  <Button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    sx={{
                      color: activeTab === tab.id ? 'primary.main' : 'text.secondary',
                      fontWeight: activeTab === tab.id ? 700 : 500,
                      '&:hover': {
                        color: 'primary.main',
                        bgcolor: 'rgba(0, 0, 0, 0.04)',
                      },
                    }}
                  >
                    {tab.label}
                  </Button>
                ))}
                <Button 
                  variant="outlined" 
                  onClick={handleLogout}
                  sx={{ ml: 1 }}
                >
                  Выйти
                </Button>
              </Box>
            )}
          </Toolbar>
        </AppBar>

        <Container maxWidth="xl" sx={{ py: 4 }}>
          {getTabContent()}
        </Container>
      </Box>
      <ToastContainer
        position="top-right"
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="light"
        style={{ marginTop: '64px' }}
      />
    </ThemeProvider>
  );
}

export default App;
