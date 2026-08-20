import {
  useState,
  type MouseEvent,
  type ReactNode,
} from "react";

import {
  AppBar,
  Avatar,
  Box,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  Toolbar,
  Tooltip,
  Typography,
} from "@mui/material";

import {
  DashboardOutlined,
  GroupOutlined,
  Logout,
  Menu as MenuIcon,
  PeopleAltOutlined,
  WarningAmberOutlined,
  AccountBalanceOutlined,
  PictureAsPdfOutlined,
} from "@mui/icons-material";

import {
  Outlet,
  useLocation,
  useNavigate,
} from "react-router-dom";

import { useAuth } from "../contexts/AuthContext";
import FGLogo from "../components/common/FGLogo";

const drawerWidth = 260;

interface NavigationItem {
  label: string;
  icon: ReactNode;
  path: string;
}

const navigationItems: NavigationItem[] = [
  {
    label: "Dashboard",
    icon: <DashboardOutlined />,
    path: "/",
  },
  {
    label: "Loan Management",
    icon: <AccountBalanceOutlined />,
    path: "/loans",
  },
  {
    label: "Client Details",
    icon: <PeopleAltOutlined />,
    path: "/clients",
  },
  {
    label: "Portfolio at Risk",
    icon: <WarningAmberOutlined />,
    path: "/portfolio-risk",
  },
  {
    label: "User Management",
    icon: <GroupOutlined />,
    path: "/users",
  },
  {
    label: "Reports",
    path: "/reports",
    icon: <PictureAsPdfOutlined />,
  },
];

export default function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();

  const {
    user,
    logout,
  } = useAuth();

  const [
    mobileOpen,
    setMobileOpen,
  ] = useState(false);

  const [
    anchorEl,
    setAnchorEl,
  ] = useState<null | HTMLElement>(null);

  const menuOpen = Boolean(anchorEl);

  const handleDrawerToggle = () => {
    setMobileOpen(
      (current) => !current,
    );
  };

  const handleProfileClick = (
    event: MouseEvent<HTMLElement>,
  ) => {
    setAnchorEl(
      event.currentTarget,
    );
  };

  const handleProfileClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleProfileClose();

    logout();

    navigate("/login", {
      replace: true,
    });
  };

  const handleNavigation = (
    path: string,
  ) => {
    navigate(path);
    setMobileOpen(false);
  };

  const isActive = (
    path: string,
  ) => {
    if (path === "/") {
      return location.pathname === "/";
    }

    return location.pathname.startsWith(
      path,
    );
  };

  const drawer = (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        backgroundColor: "#8F2115",
        color: "#FFFFFF",
      }}
    >
      {/* Logo */}
      <Box
        sx={{
          px: 2.5,
          py: 2.5,
          display: "flex",
          alignItems: "center",
          gap: 1.5,
        }}
      >
        <FGLogo
          size={52}
          sx={{
            borderRadius: 1.5,
          }}
        />

        <Box>
          <Typography
            variant="subtitle1"
            fontWeight={800}
            sx={{
              color: "#FFFFFF",
              lineHeight: 1.2,
            }}
          >
            FG
          </Typography>

          <Typography
            variant="caption"
            sx={{
              color: "#E0C477",
              lineHeight: 1.2,
            }}
          >
            Financial Friend
          </Typography>
        </Box>
      </Box>

      <Divider
        sx={{
          borderColor:
            "rgba(255,255,255,0.18)",
        }}
      />

      {/* Navigation */}
      <List
        sx={{
          px: 1.5,
          py: 2,
          flexGrow: 1,
        }}
      >
        {navigationItems.map(
          (item) => {
            const active = isActive(
              item.path,
            );

            return (
              <ListItemButton
                key={item.path}
                selected={active}
                onClick={() =>
                  handleNavigation(
                    item.path,
                  )
                }
                sx={{
                  position: "relative",
                  minHeight: 48,
                  mb: 0.75,
                  px: 1.75,
                  borderRadius: 1.5,
                  color:
                    "rgba(255,255,255,0.82)",

                  "& .MuiListItemIcon-root":
                    {
                      color:
                        "rgba(255,255,255,0.72)",
                      minWidth: 42,
                    },

                  "&:hover": {
                    backgroundColor:
                      "rgba(255,255,255,0.09)",
                    color: "#FFFFFF",

                    "& .MuiListItemIcon-root":
                      {
                        color: "#E0C477",
                      },
                  },

                  "&.Mui-selected": {
                    backgroundColor:
                      "rgba(255,255,255,0.14)",
                    color: "#FFFFFF",

                    "& .MuiListItemIcon-root":
                      {
                        color: "#E0C477",
                      },

                    "&:hover": {
                      backgroundColor:
                        "rgba(255,255,255,0.18)",
                    },
                  },

                  ...(active && {
                    "&::before": {
                      content: '""',
                      position:
                        "absolute",
                      left: 0,
                      top: 8,
                      bottom: 8,
                      width: 4,
                      borderRadius:
                        "0 4px 4px 0",
                      backgroundColor:
                        "#D0B050",
                    },
                  }),
                }}
              >
                <ListItemIcon>
                  {item.icon}
                </ListItemIcon>

                <ListItemText
                  primary={item.label}
                  primaryTypographyProps={{
                    fontSize: 14,
                    fontWeight: active
                      ? 700
                      : 500,
                  }}
                />
              </ListItemButton>
            );
          },
        )}
      </List>

      {/* Sidebar Footer */}
      <Box
        sx={{
          px: 2,
          py: 2,
          borderTop:
            "1px solid rgba(255,255,255,0.14)",
        }}
      >
        <Typography
          variant="caption"
          sx={{
            display: "block",
            color:
              "rgba(255,255,255,0.55)",
            textAlign: "center",
          }}
        >
          FG Financial Friend
        </Typography>

        <Typography
          variant="caption"
          sx={{
            display: "block",
            color: "#E0C477",
            textAlign: "center",
            mt: 0.25,
          }}
        >
          Financial Management System
        </Typography>
      </Box>
    </Box>
  );

  return (
    <Box
      sx={{
        display: "flex",
        minHeight: "100vh",
        backgroundColor:
          "#F8F5F3",
      }}
    >
      {/* Mobile Navigation */}
      <Box
        component="nav"
        sx={{
          width: {
            md: drawerWidth,
          },
          flexShrink: {
            md: 0,
          },
        }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={
            handleDrawerToggle
          }
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: {
              xs: "block",
              md: "none",
            },

            "& .MuiDrawer-paper": {
              boxSizing:
                "border-box",
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>

        {/* Desktop Navigation */}
        <Drawer
          variant="permanent"
          open
          sx={{
            display: {
              xs: "none",
              md: "block",
            },

            "& .MuiDrawer-paper": {
              boxSizing:
                "border-box",
              width: drawerWidth,
              border: "none",
            },
          }}
        >
          {drawer}
        </Drawer>
      </Box>

      {/* Main Application */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          width: {
            md: `calc(100% - ${drawerWidth}px)`,
          },
          minWidth: 0,
          minHeight: "100vh",
        }}
      >
        {/* Top Bar */}
        <AppBar
          position="fixed"
          color="inherit"
          elevation={0}
          sx={{
            width: {
              md: `calc(100% - ${drawerWidth}px)`,
            },
            ml: {
              md: `${drawerWidth}px`,
            },
            backgroundColor:
              "#FFFFFF",
            borderBottom:
              "1px solid #E7DDD9",
          }}
        >
          <Toolbar
            sx={{
              minHeight: {
                xs: 64,
                md: 70,
              },
            }}
          >
            {/* Mobile Menu */}
            <IconButton
              edge="start"
              onClick={
                handleDrawerToggle
              }
              sx={{
                mr: 2,
                display: {
                  md: "none",
                },
                color: "#8F2115",
              }}
            >
              <MenuIcon />
            </IconButton>

            {/* Page Title */}
            <Box
              sx={{
                flexGrow: 1,
              }}
            >
              <Typography
                variant="h6"
                fontWeight={700}
                sx={{
                  color: "#2B211F",
                }}
              >
                {location.pathname ===
                "/"
                  ? "Dashboard"
                  : navigationItems.find(
                      (item) =>
                        isActive(
                          item.path,
                        ),
                    )?.label ||
                    "FG Financial Friend"}
              </Typography>

              <Typography
                variant="caption"
                sx={{
                  color: "#756B68",
                }}
              >
                Financial Management System
              </Typography>
            </Box>

            {/* User */}
            <Tooltip title="Account">
              <IconButton
                onClick={
                  handleProfileClick
                }
                sx={{
                  ml: 1,
                }}
              >
                <Avatar
                  sx={{
                    width: 40,
                    height: 40,
                    backgroundColor:
                      "#8F2115",
                    color: "#FFFFFF",
                    fontWeight: 700,
                  }}
                >
                  {(
                    user?.full_name ||
                    user?.username ||
                    "U"
                  )
                    .charAt(0)
                    .toUpperCase()}
                </Avatar>
              </IconButton>
            </Tooltip>

            <Menu
              anchorEl={anchorEl}
              open={menuOpen}
              onClose={
                handleProfileClose
              }
              PaperProps={{
                sx: {
                  mt: 1,
                  minWidth: 220,
                  border:
                    "1px solid #E7DDD9",
                },
              }}
            >
              <MenuItem
                disabled
                sx={{
                  opacity: 1,
                }}
              >
                <Box>
                  <Typography
                    variant="body2"
                    fontWeight={700}
                    sx={{
                      color: "#2B211F",
                    }}
                  >
                    {user?.full_name ||
                      user?.username}
                  </Typography>

                  <Typography
                    variant="caption"
                    sx={{
                      color: "#756B68",
                    }}
                  >
                    {user?.role}
                  </Typography>
                </Box>
              </MenuItem>

              <Divider />

              <MenuItem
                onClick={
                  handleLogout
                }
                sx={{
                  color: "#8F2115",
                }}
              >
                <ListItemIcon
                  sx={{
                    minWidth: 36,
                    color: "#8F2115",
                  }}
                >
                  <Logout
                    fontSize="small"
                  />
                </ListItemIcon>

                <ListItemText>
                  Logout
                </ListItemText>
              </MenuItem>
            </Menu>
          </Toolbar>
        </AppBar>

        {/* Page Content */}
        <Box
          sx={{
            pt: {
              xs: 10,
              md: 11,
            },
            px: {
              xs: 2,
              sm: 3,
              md: 4,
            },
            pb: 4,
          }}
        >
          <Outlet />
        </Box>
      </Box>
    </Box>
  );
}