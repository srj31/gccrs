/* Configuration for GNU C-compiler for MIPS Rx000 family
   Copyright (C) 1989, 1990, 1991, 1993, 1997, 2001
   Free Software Foundation, Inc.

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
#if !defined(MIPSEL) && !defined(__MIPSEL__)
#define HOST_WORDS_BIG_ENDIAN
#endif

/* A code distinguishing the floating point format of the host
   machine.  There are three defined values: IEEE_FLOAT_FORMAT,
   VAX_FLOAT_FORMAT, and UNKNOWN_FLOAT_FORMAT.  */

#define HOST_FLOAT_FORMAT IEEE_FLOAT_FORMAT

#ifndef __GNUC__
/* The MIPS compiler gets it wrong, and treats enumerated bitfields
   as signed quantities, making it impossible to use an 8-bit enum
   for compiling GNU C++.  */
#define ONLY_INT_FIELDS 1
#endif
