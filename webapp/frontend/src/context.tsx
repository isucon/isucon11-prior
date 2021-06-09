import React, { createContext, FC, useCallback, useContext, useEffect, useState } from "react";
import { User } from "./model";


interface AppContext {
  user?: User;
  login: (user: User) => void;
  logout: () => void;
}

const AppContext = createContext<AppContext>({
  user: undefined,
  login: () => {},
  logout: () => {},
});

export const useAppContext = () => useContext(AppContext);
export const AppContextProvider: FC = ({ children }) => {
  const [user, update] = useState<User | undefined>(undefined)

  const login = useCallback<AppContext['login']>((user) => {
    update(user)
  }, [update])
  const logout = useCallback<AppContext['logout']>(() => {
    update(undefined);
  }, [update])

  useEffect(() => {
    fetch('/api/session').then((r) => r.json()).then((user) => {
      update(user)
    })
  }, [update])

  return <AppContext.Provider value={{ user, login, logout }} children={children} />
}
