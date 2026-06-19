"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Zap } from "lucide-react";
import { getBalance } from "@/services/api/credits";
import { getStoredToken } from "@/services/api/client";
import { useCreditBalanceRefreshSignal } from "@/constant/credits";

export function UserCreditDisplay() {
  const [balance, setBalance] = useState<number | null>(null);
  const token = getStoredToken();
  const refreshSignal = useCreditBalanceRefreshSignal();

  useEffect(() => {
    if (!token) return;
    getBalance().then((data) => setBalance(data.balance)).catch(() => {});
  }, [refreshSignal, token]);

  if (!token || balance === null) return null;

  return (
    <Link href="/credits" className="inline-flex items-center gap-1 text-xs text-stone-500 hover:text-amber-600 dark:text-stone-400 dark:hover:text-amber-400 transition-colors" title="积分明细">
      <Zap className="size-3 fill-amber-400 text-amber-400" />
      <span>{balance.toLocaleString()}</span>
    </Link>
  );
}
