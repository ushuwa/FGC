import {
  Navigate,
  Route,
  Routes,
} from "react-router-dom";

import LoginPage from "../pages/auth/LoginPage";
import DashboardPage from "../pages/dashboard/DashboardPage";

import MainLayout from "../layouts/MainLayout";
import ProtectedRoute from "./ProtectedRoute";

import ClientDetailsPage from "../pages/clients/ClientDetailsPage";
import ClientProfilePage from "../pages/clients/ClientProfilePage";
import LoanManagementPage from "../pages/loans/LoanManagementPage";
import LoanProfilePage from "../pages/loans/LoanProfilePage";
import PortfolioPage from "../pages/portfolio/PortfolioAtRiskPage";
import ReportPage from "../pages/reports/ReportsPage";

export default function AppRoutes() {
  return (
    <Routes>

      {/* =========================================
          PUBLIC ROUTES
      ========================================= */}

      <Route
        path="/login"
        element={
          <LoginPage />
        }
      />

      {/* =========================================
          PROTECTED ROUTES
      ========================================= */}

      <Route
        element={
          <ProtectedRoute />
        }
      >
        <Route
          element={
            <MainLayout />
          }
        >

          {/* Dashboard */}
          <Route
            index
            element={
              <DashboardPage />
            }
          />

          {/* Client List */}
          <Route
            path="clients"
            element={
              <ClientDetailsPage />
            }
          />

          {/* Client Profile */}
          <Route
            path="clients/:id"
            element={
              <ClientProfilePage />
            }
          />

          {/* Loan Management */}
          <Route
            path="loans"
            element={
              <LoanManagementPage />
            }
          />

          {/* Loan Profile */}
          <Route
            path="loans/:id"
            element={
              <LoanProfilePage />
            }
          />

          {/* Portfolio */}
          <Route
            path="portfolio-risk"
            element={
              <PortfolioPage />
            }
          />

          {/* Reports */}
          <Route
            path="reports"
            element={
              <ReportPage />
            }
          />

        </Route>
      </Route>

      {/* =========================================
          FALLBACK
      ========================================= */}

      <Route
        path="*"
        element={
          <Navigate
            to="/"
            replace
          />
        }
      />

    </Routes>
  );
}