"use client";

import { useEffect, useState } from "react";
import { Zap } from "lucide-react";
import { getBalance } from "@/services/api/credits";
import { getStoredToken } from "@/services/api/client";

export function UserCreditDisplay() {
  const [balance, setBalance] = useState<number | null>(null);
  const token = getStoredToken();

  useEffect(() => {
    if (!token) return;
    getBalance().then((data) => setBalance(data.balance)).catch(() => {});
  }, [token]);

  if (!token || balance === null) return null;

  return (
    <span className="inline-flex items-center gap-1 text-xs text-stone-500 dark:text-stone-400">
      <Zap className="size-3 fill-amber-400 text-amber-400" />
      <span>{balance.toLocaleString()}</span>
    </span>
  );
}
