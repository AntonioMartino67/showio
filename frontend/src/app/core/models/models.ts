export interface User {
  id: string;
  username: string;
  email: string;
  avatar_url?: string;
  has_password?: boolean;
  google_linked?: boolean;
  notify_new_seasons?: boolean;
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
  tags: Tag[];
}

export interface UpcomingEpisode {
  media_item_id: string;
  title: string;
  poster_url?: string;
  season_number: number;
  episode_number: number;
  air_date?: string;
}

export interface EpisodeDetail {
  season_number: number;
  episode_number: number;
  title?: string;
  air_date?: string;
}

export interface Tag {
  id: string;
  name: string;
  color: string;
}

export interface MediaDetail {
  media_item_id: string;
  title: string;
  type: string;
  poster_url?: string;
  overview?: string;
  status?: ProgressStatus;
  current_season: number;
  current_episode: number;
  rating?: number;
  episodes: EpisodeDetail[];
  tags: Tag[];
}

export interface Stats {
  total_titles: number;
  by_status: Record<string, number>;
  by_type: Record<string, number>;
  total_episodes_watched: number;
  average_rating?: number;
  rated_titles_count: number;
}