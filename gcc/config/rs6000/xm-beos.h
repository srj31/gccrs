/* Configuration for GNU C-compiler for BeOS host.
   Copyright (C) 1997, 1999, 2001 Free Software Foundation, Inc.
   Contributed by Fred Fish (fnf@cygnus.com), based on xm-rs6000.h
   by Richard Kenner (kenner@vlsi1.ultra.nyu.edu).


This file is part of GNU CC.

GNU CC is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2, or (at your option)
any later version.

GNU CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GNU CC; see the file COPYING.  If not, write to
the Free Software Foundation, 59 Temple Place - Suite 330,
Boston, MA 02111-1307, USA.  */

/* This describes the machine the compiler is hosted on.  */
#define	HOST_WORDS_BIG_ENDIAN

/* Use only int bitfields.  */
#define	ONLY_INT_FIELDS

/* use ANSI/SYSV style byte manipulation routines instead of BSD ones */

#undef bcopy
#define bcopy(s,d,n)	memmove((d),(s),(n))

/* BeOS is closer to USG than BSD */

#define USG

/* Define various things that the BeOS host has. */

#ifndef HAVE_VPRINTF
#define HAVE_VPRINTF
#endif
#ifndef HAVE_PUTENV
#define HAVE_PUTENV
#endif
#ifndef HAVE_RENAME
#define HAVE_RENAME
#endif

/* STANDARD_INCLUDE_DIR is the equivalent of "/usr/include" on UNIX. */

#define STANDARD_INCLUDE_DIR	"/boot/develop/headers/posix"

/* SYSTEM_INCLUDE_DIR is the location for system specific, non-POSIX headers. */

#define SYSTEM_INCLUDE_DIR	"/boot/develop/headers/be"

/* This is a temporary hack until the wimpy default 64k stack
   limit in BeOS is either increased or made user settable somehow.
   This probably won't happen until after the DR9 release.  */
#undef USE_C_ALLOCA
#define USE_C_ALLOCA 1
