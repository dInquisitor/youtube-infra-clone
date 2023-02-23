import { useState } from "react";

const useLocalStorage = (keyName, defaultValue) => {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const value = window.localStorage.getItem(keyName);
      if (value) {
        return JSON.parse(value);
      }
      window.localStorage.setItem(keyName, JSON.stringify(defaultValue));
      return defaultValue;
    } catch (err) {
      return defaultValue;
    }
  });
  const setValue = (newValue) => {
    window.localStorage.setItem(keyName, JSON.stringify(newValue));
    setStoredValue(newValue);
  };
  return [storedValue, setValue];
};

export default useLocalStorage;
