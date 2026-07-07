export interface User {
  id: string;
  username: string;
  email: string;
}

export interface AuthResponse {
  token: string;
}

export interface SearchResult {
  media_item_id: string;
  external_id: string;
  source: 'tmdb' | 'anilist';
  title: string;
  type: 'movie' | 'tv' | 'anime';
  poster_url: string;
  overview: string;
}

export type ProgressStatus = 'watching' | 'completed' | 'dropped' | 'plan_to_watch';

export interface ProgressItem {
  progress_id: string;
  status: ProgressStatus;
  current_season: number;
  current_episode: number;
  rating?: number;
  media_item_id: string;
  title: string;
  type: string;
  poster_url?: string;
}

export interface UpcomingEpisode {
  media_item_id: string;
  title: string;
  poster_url?: string;
  season_number: number;
  episode_number: number;
  air_date?: string;
}