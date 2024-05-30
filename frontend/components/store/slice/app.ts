import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";

interface App {
  id: string;
  name: string;
  namespace: string;
  deployedAt: number;
  status: string;
}

interface AppState {
  apps: App[];
}

const initialState: AppState = {
  apps: [],
};

export const fetchApps = createAsyncThunk("app/fetchApps", async () => {
  const response = await axios.get("/app");
  return response.data;
});

export const deleteApp = createAsyncThunk(
  "app/deleteApp",
  async (id: string) => {
    await axios.delete(`/apps/${id}`);
    return id;
  }
);

const appSlice = createSlice({
  name: "app",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchApps.fulfilled, (state, action) => {
      state.apps = action.payload;
    });
    builder.addCase(deleteApp.fulfilled, (state, action) => {
      state.apps = state.apps.filter((app) => app.id !== action.payload);
    });
  },
});

export default appSlice.reducer;
